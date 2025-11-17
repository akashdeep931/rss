package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	Item        []RSSItem `xml:"item"`
}

type RSSFeed struct {
	Channel RSSChannel `xml:"channel"`
}

func urlToFeed(url string) (rssFeed RSSFeed, err error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := httpClient.Get(url)
	if err != nil {
		return rssFeed, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return rssFeed, err
	}

	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return rssFeed, err
	}

	return rssFeed, nil
}
