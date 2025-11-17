package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/akashdeep931/rss/internal/db"
	"github.com/google/uuid"
)

func startScraping(db *db.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v goroutines every %s", concurrency, timeBetweenRequest)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(ctx, int32(concurrency))
		if err != nil {
			log.Println("error fetching feed:", err.Error())
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(wg, ctx, db, feed)
		}

		wg.Wait()
	}
}

func scrapeFeed(wg *sync.WaitGroup, ctx context.Context, dbConn *db.Queries, feed db.Feed) {
	defer wg.Done()

	_, err := dbConn.MarkFeedAsFetched(ctx, feed.ID)
	if err != nil {
		log.Println("error marking feed as fetched:", err.Error())
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("error fetching feed:", err.Error())
		return
	}

	for _, item := range rssFeed.Channel.Item {
		var description sql.NullString

		if len(item.Description) > 0 {
			description = sql.NullString{String: item.Description, Valid: true}
		}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("couldn't parse date %v with err %s", item.PubDate, err.Error())
			continue
		}

		_, err = dbConn.CreatePost(ctx, db.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}

			log.Println("failed to create post:", err.Error())
		}
	}

	log.Printf("Feed %s collected, %d posts found", feed.Name, len(rssFeed.Channel.Item))
}
