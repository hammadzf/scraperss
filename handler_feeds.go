package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/hammadzf/scraperss/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON in the request body: %v", err))
	}

	// check if a feed with the same ULR already exists
	_, err = apiCfg.DB.GetFeedByURL(r.Context(), database.GetFeedByURLParams{
		UserID: user.ID,
		Url:    params.URL,
	})
	if err == nil {
		respondWithError(w, 400, "An RSS feed with this URL already exists.")
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		Name:      params.Name,
		Url:       params.URL,
		ID:        uuid.New(),
		UserID:    user.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't create feed: %v", err))
	}
	respondWithJSON(w, 201, databaseFeedToFeed(feed))
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := apiCfg.DB.GetFeedsOfUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error fetching feeds: %v", err))
		return
	}
	if feeds == nil {
		respondWithError(w, 404, fmt.Sprintf("No feeds exist for user with ID %v", user.ID))
		return
	}
	respondWithJSON(w, 200, databaseFeedsToFeeds(feeds))
}

func (apiCfg *apiConfig) handlerDeleteFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	feedIdStr := chi.URLParam(r, "feedID")
	feedId, err := uuid.Parse(feedIdStr)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error parsing feed ID: %v", err))
		return
	}
	err = apiCfg.DB.DeleteFeed(r.Context(), database.DeleteFeedParams{
		UserID: user.ID,
		ID:     feedId,
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't delete feed: %v", err))
		return
	}
	respondWithJSON(w, 204, struct{}{})
}
