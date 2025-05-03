package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/hammadzf/scraperss/internal/database"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// request format for POST /users operation
	type parameters struct {
		Name string `json:"name"`
	}

	// decode request body as per format
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf(
			"Error parsing JSON in the request body: %v", err),
		)
		return
	}

	usr, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error creating user: %v", err))
		return
	}
	respondWithJSON(w, 201, databaseUserToUser(usr))
}

func (apiCfg *apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	usrs, err := apiCfg.DB.GetUsers(r.Context())
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error fetching users: %v", err))
		return
	}
	if usrs == nil {
		respondWithError(w, 404, "No users found.")
		return
	}
	respondWithJSON(w, 200, databaseUsersToUsers(usrs))
}

func (apiCfg *apiConfig) handlerGetUserById(w http.ResponseWriter, r *http.Request) {
	usrIdStr := chi.URLParam(r, "userID")
	usrId, err := uuid.Parse(usrIdStr)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error parsing user ID: %v", err))
		return
	}
	usr, err := apiCfg.DB.GetUserByID(r.Context(), usrId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, 404,
				fmt.Sprintf("User with ID %v does not exist.", usrId))
			return
		}
		respondWithError(w, 500, fmt.Sprintf("Error fetching users: %v", err))
		return
	}
	respondWithJSON(w, 200, databaseUserToUser(usr))
}

func (apiCfg *apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	usrIdStr := chi.URLParam(r, "userID")
	usrId, err := uuid.Parse(usrIdStr)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error parsing user ID: %v", err))
		return
	}
	err = apiCfg.DB.DeleteUser(r.Context(), usrId)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete user: %v", err))
		return
	}
	respondWithJSON(w, 204, struct{}{})
}
