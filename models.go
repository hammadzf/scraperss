package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/hammadzf/scraperss/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"apiKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Feed struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Url           string    `json:"url"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	UserID        uuid.UUID `json:"userId"`
	LastFetchedAt time.Time `json:"lastFetchedAt"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		ApiKey:    dbUser.ApiKey,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}

func databaseFeedToFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:            dbFeed.ID,
		Name:          dbFeed.Name,
		Url:           dbFeed.Url,
		CreatedAt:     dbFeed.CreatedAt,
		UpdatedAt:     dbFeed.UpdatedAt,
		UserID:        dbFeed.UserID,
		LastFetchedAt: dbFeed.LastFetchedAt.Time,
	}
}

func databaseUsersToUsers(dbUsers []database.User) []User {
	users := []User{}
	for _, dbUser := range dbUsers {
		users = append(users, databaseUserToUser(dbUser))
	}
	return users
}

func databaseFeedsToFeeds(dbFeeds []database.Feed) []Feed {
	feeds := []Feed{}
	for _, dbFeed := range dbFeeds {
		feeds = append(feeds, databaseFeedToFeed(dbFeed))
	}
	return feeds
}
