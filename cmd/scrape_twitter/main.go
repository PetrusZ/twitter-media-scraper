package main

import (
	"context"
	"flag"
	"fmt"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

var  (
	twitterUser = flag.String("user", "", "the twitter user to get files from")
	tweetAmount = flag.Int("amount", 1000, "amount of tweets to get content from")
	getVideos = flag.Bool("videos", true, "download videos from tweets")
	getPhotos = flag.Bool("photos", true, "download photos from tweets")
)

func main() {

	if *twitterUser != "" {
		GetUserTweets(*twitterUser, *tweetAmount)
	}
}

func GetUserTweets(user string, amount int) {
	scraper := twitterscraper.New()

	downloader := FileDownloader{FileUrls: make(chan string)}
	downloader.Start(15)
	for tweet := range scraper.GetTweets(context.Background(), user, amount) {
		fmt.Println("found tweet, id: ", tweet.ID)
		if tweet.Error != nil {
			panic(tweet.Error)
		}

		if *getVideos {
			for _, video := range tweet.Videos {
				fmt.Println("downloading: ", video.URL)
				downloader.FileUrls <- video.URL
			}
		}

		if *getPhotos {
			for _, photo := range tweet.Photos {
				downloader.FileUrls <- photo
			}
		}
	}
	downloader.Done()
}

