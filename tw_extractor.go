package main

import (
	"fmt"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func createClient() *twitter.Client {
	// oauth2 configures a client that uses app credentials to keep a fresh token
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("TWITTER_CONSUMER_KEY"),
		ClientSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth2.NoContext)

	// Twitter client
	client := twitter.NewClient(httpClient)

	return client
}

func getTweetsFromPastMonths(username string, client *twitter.Client, c chan []twitter.Tweet) {
	var lastID int64
	for i := 0; i < 5; i++ {
		timeline, _, _ := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
			ScreenName: username, MaxID: lastID})

		lastID = timeline[len(timeline)-1].ID

		c <- timeline
	}
	close(c)
}

func main() {
	client := createClient()

	c := make(chan []twitter.Tweet, 5)
	go getTweetsFromPastMonths("viyuelaeveryday", client, c)
	for tl := range c {
		for _, tweet := range tl {
			fmt.Println(tweet.FavoriteCount, tweet.RetweetCount, tweet.ReplyCount, tweet.QuoteCount, tweet.CreatedAt)
		}
		fmt.Print("-------------------------------------------------------\n")
	}
}
