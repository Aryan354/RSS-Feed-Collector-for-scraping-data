package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/Aryan354.RssServer/internal/database"
	"github.com/google/uuid"
)

// startScraping initiates the perpetual scraping process.
func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequests time.Duration,
) {
	log.Printf("The perpetual scraper is running")
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for range ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Print("Error fetching feeds:", err)
			continue
		}

		var wg sync.WaitGroup
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(db, &wg, feed)
		}
		wg.Wait()
	}
}

// scrapeFeed fetches and processes the RSS feed, then stores the items in the database.
func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	if _, err := db.MarkFeedAsFetched(context.Background(), feed.ID); err != nil {
		log.Print("There was a concurrency issue:", err)
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching URL:", err)
		return
	}

	for _, item := range rssFeed.Channel.Items {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Println("Error parsing date:", err)
			continue
		}

		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: publishedAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		}

		if _, err := db.CreatePost(context.Background(), postParams); err != nil {
			log.Println("Failed to create post:", err)
		}
	}
}
