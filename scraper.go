package main

import (
	"context"
	"database/sql"
	"goapi/internal/database"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequests time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("Error fetching feed: ", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			// increment wait group counter by 1
			wg.Add(1)
			go scrapeFeed(wg, db, &feed)
		}
		// blocks till the wait group counter is zero
		wg.Wait()
	}
}

func scrapeFeed(wg *sync.WaitGroup, db *database.Queries, feed *database.Feed) {
	// decrements the wait group counter by 1
	defer wg.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error making feed as fetched: ", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed: ", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {

		desc := sql.NullString{}
		// desc.Valid is set to false by default
		if item.Description != "" {
			desc.Valid = true
			desc.String = item.Description
		}
		date, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Couldn't parse date %v with error %v", date, err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: desc,
			PublishedAt: date,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Print("Failed to create post with error: ", err)
		}
	}
	log.Printf("Feeds %s collected, %d posts found", feed.Name, len(rssFeed.Channel.Item))
}
