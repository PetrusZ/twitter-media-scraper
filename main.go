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
	var configPath string
	flag.StringVar(&configPath, "configPath", ".", "Input config file path")
	flag.Parse()

	config, err := LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	logLevel := zerolog.InfoLevel
	switch *config.Global.LogLevel {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "fatal":
		logLevel = zerolog.FatalLevel
	case "panic":
		logLevel = zerolog.PanicLevel
	case "no":
		logLevel = zerolog.NoLevel
	case "disabled":
		logLevel = zerolog.Disabled
	}

	zerolog.SetGlobalLevel(logLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	d := GetDownloaderInstance(16)

	log.Info().Msg("Downloader starts")

	for _, user := range config.Users {
		if *user.UserName != "" {
			err := getUserTweets(*user.UserName, *user.TweetAmount, *user.GetVideos, *user.GetPhotos, d)
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
