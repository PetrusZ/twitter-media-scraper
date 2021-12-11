package main

import (
	"context"
	"fmt"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

func main() {
	configFile, err := NewConfigFile("config.json")
	if err != nil {
		panic(err)
	}

	err = configFile.Load()
	if err != nil {
		panic(err)
	}

	d := GetDownloaderInstance(16)

	for _, config := range configFile.GetConfigs() {
		if config.UserName != "" {
			getUserTweets(config.UserName, config.TweetAmount, config.GetVideos, config.GetPhotos, d)
		} else {
			fmt.Println("No twitter user")
		}
	}

	close(d.GetInfo())
	d.Wait()
}

func getUserTweets(user string, amount int, getVideos bool, getPhotos bool, d Downloader) (err error) {
	scraper := twitterscraper.New()
	tweets := scraper.GetTweets(context.Background(), user, amount)

	for tweet := range tweets {
		if tweet.Error != nil {
			return tweet.Error
		}

		// url := "https://twitter.com/" + user + "/status/" + tweet.ID
		date := fmt.Sprintf("%d%02d%02d", tweet.TimeParsed.Year(), tweet.TimeParsed.Month(), tweet.TimeParsed.Day())

		if getVideos && tweet.Videos != nil {
			for _, video := range tweet.Videos {
				tweetInfo := tweetInfo{user, date + " - " + tweet.ID, video.URL}
				d.GetInfo() <- tweetInfo
			}
		}

		if getPhotos && tweet.Videos == nil {
			for id, url := range tweet.Photos {
				tweetInfo := tweetInfo{user, date + " - " + tweet.ID + "-" + fmt.Sprint(id), url + "?format=jpg&name=orig"}
				d.GetInfo() <- tweetInfo
			}
		}
	}

	return err
}
