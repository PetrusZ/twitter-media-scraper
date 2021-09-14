package main

import (
	"context"
	"fmt"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

var (
	twitterUser string
	tweetAmount int
	getVideos   bool
	getPhotos   bool
)

func main() {
	configFile := NewConfigFile()
	err := configFile.Load("config.json")
	if err != nil {
		panic(err)
	}

	d := GetDownloaderInstance()

	for _, config := range configFile.Configs {
		twitterUser = config.UserName
		tweetAmount = config.TweetAmount
		getVideos = config.GetVideos
		getPhotos = config.GetPhotos

		if twitterUser != "" {
			getUserTweets(twitterUser, tweetAmount, d)
		} else {
			fmt.Println("No twitter user")
		}
	}

	close(d.info)
	d.Wait()
}

func getUserTweets(user string, amount int, d *downloader) (err error) {
	scraper := twitterscraper.New()
	tweets := scraper.GetTweets(context.Background(), user, amount)

	if tweets == nil {
		err = mkdirAll("out/" + user + "/")
		if err != nil {
			return err
		}
	}

	for tweet := range tweets {
		if tweet.Error != nil {
			return tweet.Error
		}

		url := "https://twitter.com/" + user + "/status/" + tweet.ID

		if getVideos {
			if tweet.Videos != nil {
				tweetInfo := tweetInfo{user, "", url, Video}
				d.info <- tweetInfo
			}
		}

		if getPhotos {
			if tweet.Videos == nil {
				for id, url := range tweet.Photos {
					date := fmt.Sprintf("%d%02d%02d", tweet.TimeParsed.Year(), tweet.TimeParsed.Month(), tweet.TimeParsed.Day())
					tweetInfo := tweetInfo{user, date + " - " + tweet.ID + "-" + fmt.Sprint(id), url, Photo}
					d.info <- tweetInfo
				}
			}
		}
	}

	return err
}
