package main

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hammadzf/scraperss/internal/database"
)

func startScraping(db *database.Queries, concurrency int, interval time.Duration) {
	log.Printf("Scraping feeds using %v goroutines every %v duration", concurrency, interval)
	// start a time ticker
	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		// get next feeds to fetch
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("couldn't fetch feeds: %v", err)
			continue
		}
		// start go routines to scrape feeds in parallel
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	// fetch feed and mark feed as fetched
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking the feed as fetched: %v", err)
		return
	}
	// fetch feed from url
	rssFeed, err := fetchFeedFromUrl(feed.Url)
	if err != nil {
		log.Printf("couldn't fetch feed from its url: %v", err)
		return
	}

	// parse through all items on the RSS channel
	// and save them as individual posts in DB
	for _, item := range rssFeed.Channel.Item {
		pubAt, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			log.Printf("Error parsing pubDate %v with error: %v", item.PubDate, err)
			continue
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				// post already exists in the DB
				continue
			}
			log.Printf("Couldn't create post: %v", err)
		}
		log.Printf("Found post %s on feed %s", item.Title, feed.Name)
	}
	log.Printf("Collected %v posts from feed %s", len(rssFeed.Channel.Item), feed.Name)

}
