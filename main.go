package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/hammadzf/scraperss/internal/database"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// for connection to DB
type apiConfig struct {
	DB *database.Queries
}

//go:embed sql/schema/*.sql
var embedMigrations embed.FS

func main() {

	// Get DB password
	bin, err := os.ReadFile("/run/secrets/db-password")
	if err != nil {
		log.Fatal("Couldn't get DB password:", err)
	}

	// Connect to DB
	conn, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres:%s@db:5432/scraperss?sslmode=disable", string(bin)))
	if err != nil {
		log.Fatal("Couldn't connect to DB:", err)
	}

	// run goose migrations
	goose.SetBaseFS(embedMigrations)

	err = goose.SetDialect("postgres")
	if err != nil {
		log.Fatal("Couldn't set dialect for goose:", err)
	}

	err = goose.Up(conn, "sql/schema")
	if err != nil {
		log.Fatal("Couldn't run goose migrations:", err)
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
		Handler: router,
	}

	// run service and catch error
	log.Printf("Starting server on port 80")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Couldn't start server:", err)
	}
}
