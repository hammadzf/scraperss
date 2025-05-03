package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/hammadzf/scraperss/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// for connection to DB
type apiConfig struct {
	DB *database.Queries
}

func main() {
	// Load environment file
	godotenv.Load(".env")
	// Get port number
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Couldn't find PORT value in environment file.")
	}
	// Get DB URL
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Couldn't find DB URL in the environment file.")
	}

	// Connect to DB
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Couldn't connect to DB:", err)
	}
	db := database.New(conn)

	// DB Config
	apiCfg := apiConfig{
		DB: db,
	}

	// start scraping 10 feeds in parallel every 1 minute
	go startScraping(db, 10, time.Minute)

	// create router
	router := chi.NewRouter()
	// CORS configurations
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"POST", "GET", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// create internal v1 router
	v1Router := chi.NewRouter()

	// basic check endpoints
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)

	// users endpoints (unauthorized)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.handlerGetUsers)
	v1Router.Get("/users/{userID}", apiCfg.handlerGetUserById)
	v1Router.Delete("/users/{userID}", apiCfg.handlerDeleteUser)

	// feeds endpoints (authorized)
	v1Router.Post("/feeds", apiCfg.middlewareAuthzHandler(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.middlewareAuthzHandler(apiCfg.handlerGetFeeds))
	v1Router.Delete("/feeds/{feedID}", apiCfg.middlewareAuthzHandler(apiCfg.handlerDeleteFeed))

	// mount v1 router to the main router
	router.Mount("/v1", v1Router)

	// create server
	srv := http.Server{
		Addr:    ":" + portString,
		Handler: router,
	}

	// run service and catch error
	fmt.Printf("Starting server on port %s", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Couldn't start server:", err)
	}
}
