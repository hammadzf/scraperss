package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title string    `xml:"title"`
		Link  string    `xml:"link"`
		Item  []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
}

func fetchFeedFromUrl(url string) (RSSFeed, error) {
	// HTTP client to fetch the feed from the URL
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}

	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSSFeed{}, err
	}

	rssFeed := RSSFeed{}
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return RSSFeed{}, err
	}

	return rssFeed, nil
}
