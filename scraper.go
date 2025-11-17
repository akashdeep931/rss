package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/akashdeep931/rss/internal/db"
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

func scrapeFeed(wg *sync.WaitGroup, ctx context.Context, db *db.Queries, feed db.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(ctx, feed.ID)
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
		log.Println("Found post", item.Title, "on feed", feed.Name)
	}
	log.Printf("Feed %s collected, %d posts found", feed.Name, len(rssFeed.Channel.Item))
}
