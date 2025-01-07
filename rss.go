package main

import (
	"encoding/xml"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func urlToFeed(url string) (*RSSFeed, error) {
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	res, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// dat, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	return nil, err
	// }

	// var rssFeed RSSFeed
	// err = xml.Unmarshal(dat, &rssFeed)
	// if err != nil {
	// 	return nil, err
	// }

	var rssFeed RSSFeed
	decoder := xml.NewDecoder(res.Body)
	err = decoder.Decode(&rssFeed)
	if err != nil {
		return nil, err
	}

	return &rssFeed, nil
}
