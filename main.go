package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var configName string
	flag.StringVar(&configName, "configFile", "config.json", "Input configFile name")
	flag.Parse()

	configFile := NewConfigFile(configName)

	err := configFile.Load()
	if err != nil {
		panic(err)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	d := GetDownloaderInstance(16)

	log.Info().Msg("Downloader starts")

	for _, config := range configFile.GetConfigs() {
		if config.UserName != "" {
			err := getUserTweets(config.UserName, config.TweetAmount, config.GetVideos, config.GetPhotos, d)
			if err != nil {
				log.Error().Err(err).Msg("")
			}
		} else {
			log.Error().Msg("No twitter user found")
		}
	}

	close(d.GetInfo())
	d.Wait()
	d.PrintCounter()
}

func getUserTweets(user string, amount int, getVideos bool, getPhotos bool, d Downloader) (err error) {
	log.Debug().Msgf("Downloading user %s's video = %t, photos = %t", user, getVideos, getPhotos)

	scraper := twitterscraper.New()
	// scraper.WithDelay(60)
	tweets := scraper.GetTweets(context.Background(), user, amount)

	for tweet := range tweets {
		if tweet.Error != nil {
			return fmt.Errorf("%s tweet.Error: %w", user, tweet.Error)
		}

		// url := "https://twitter.com/" + user + "/status/" + tweet.ID
		date := fmt.Sprintf("%d%02d%02d", tweet.TimeParsed.Year(), tweet.TimeParsed.Month(), tweet.TimeParsed.Day())

		if getVideos && tweet.Videos != nil {
			for _, video := range tweet.Videos {
				tweetInfo := tweetInfo{user, user, date + " - " + tweet.ID, video.URL, TweetTypeVideo}
				d.GetInfo() <- tweetInfo
			}
		}

		if getPhotos && tweet.Videos == nil {
			for id, url := range tweet.Photos {
				tweetInfo := tweetInfo{user, user, date + " - " + tweet.ID + "-" + fmt.Sprint(id), url + "?format=jpg&name=orig", TweetTypePhoto}
				d.GetInfo() <- tweetInfo
			}
		}
	}

	return err
}
