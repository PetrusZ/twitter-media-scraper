package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"syscall"

	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/PetrusZ/twitter-media-scraper/internal/config"
	"github.com/PetrusZ/twitter-media-scraper/internal/downloader"
	"github.com/PetrusZ/twitter-media-scraper/internal/utils"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "configPath", "./configs", "Input config file path")
	flag.Parse()

	conf, err := config.Load(configPath)
	if err != nil {
		panic(err)
	}

	config.Watch()

	logLevel := zerolog.InfoLevel
	switch *conf.Global.LogLevel {
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

	utils.Sigs = make(chan os.Signal)
	d := downloader.GetDownloaderInstance(16)

	log.Info().Msg("Downloader starts")

	for {
		log.Info().Msg("download start")

		for _, user := range conf.Users {
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

		log.Info().Msg("download end")

	SIGNAL:
		for {
			select {
			case sig := <-utils.Sigs:
				if sig == syscall.SIGHUP {
					conf, err = config.Load(configPath)
					log.Info().Msg("config reloaded")
					if err != nil {
						log.Error().Err(err).Msg("reload config error")
					} else {
						break SIGNAL
					}
				}
			}
		}

		d.Init()
	}
}

func getUserTweets(user string, amount int, getVideos bool, getPhotos bool, d downloader.Downloader) (err error) {
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
				tweetInfo := downloader.TweetInfo{
					User:      user,
					Dir:       user,
					Name:      date + " - " + tweet.ID,
					URL:       video.URL,
					TweetType: downloader.TweetTypeVideo,
				}
				d.GetInfo() <- tweetInfo
			}
		}

		if getPhotos && tweet.Videos == nil {
			for id, url := range tweet.Photos {
				tweetInfo := downloader.TweetInfo{
					User:      user,
					Dir:       user,
					Name:      date + " - " + tweet.ID + "-" + fmt.Sprint(id),
					URL:       url + "?format=jpg&name=orig",
					TweetType: downloader.TweetTypePhoto,
				}
				d.GetInfo() <- tweetInfo
			}
		}
	}

	return err
}
