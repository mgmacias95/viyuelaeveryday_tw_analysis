package main

import (
	"os"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
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
		timeline, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
			ScreenName: username, MaxID: lastID})
		if err != nil {
			panic(err)
		}

		lastID = timeline[len(timeline)-1].ID

		c <- timeline
	}
	close(c)
}

func getPlotData(client *twitter.Client) ([]int, []int, []string) {
	c := make(chan []twitter.Tweet, 5)
	var dates []string
	var favs, rts []int
	go getTweetsFromPastMonths("viyuelaeveryday", client, c)
	for tl := range c {
		for _, tweet := range tl {
			dates = append(dates, tweet.CreatedAt)
			favs = append(favs, tweet.FavoriteCount)
			rts = append(rts, tweet.RetweetCount)
		}
	}
	return favs, rts, dates
}

func generatePoints(vals []int) plotter.XYs {
	pts := make(plotter.XYs, len(vals))
	for i := range pts {
		pts[i].X = float64(i)
		pts[i].Y = float64(vals[i])
	}
	return pts
}

type dateTicks struct {
	dates []string
}

func (d dateTicks) Ticks(min, max float64) []plot.Tick {
	tks := plot.DefaultTicks{}.Ticks(min, max)
	for i, t := range tks {
		if t.Label == "" { // Skip minor ticks, they are fine.
			continue
		}
		val, _ := strconv.Atoi(t.Label)
		tks[i].Label = d.dates[val]
	}
	return tks
}

func plotData(rts []int, favs []int, dates []string) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Interactions with viyuelaeveryday"
	p.X.Label.Text = "Day"
	p.Y.Label.Text = "Number of interactions"
	p.X.Tick.Marker = dateTicks{dates}

	err = plotutil.AddLinePoints(p,
		"Favs", generatePoints(favs),
		"RTs", generatePoints(rts))
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(20*vg.Inch, 10*vg.Inch, "interactions.png"); err != nil {
		panic(err)
	}
}

func main() {
	client := createClient()
	favs, rts, dates := getPlotData(client)
	plotData(rts, favs, dates)
}
