package main

import (
	"context"
	"fmt"
	"os"
	"syscall"

	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/PetrusZ/twitter-media-scraper/internal/config"
	"github.com/PetrusZ/twitter-media-scraper/internal/downloader"
	"github.com/PetrusZ/twitter-media-scraper/internal/utils"
)

func main() {
	flag.String("log_level", "info", "log level")
	flag.String("config_path", "./configs", "config file path")
	flag.Bool("keep_running", true, "whether keep running all the time")
	flag.Parse()

	viper.BindPFlags(flag.CommandLine)
	configPath := viper.GetString("config_path")

	conf, err := config.Load(configPath)
	if err != nil {
		panic(err)
	}

	keepRunning := viper.GetBool("keep_running")
	logLevel := viper.GetString("log_level")

	config.Watch()

	zeroLogLevel := zerolog.InfoLevel
	switch logLevel {
	case "debug":
		zeroLogLevel = zerolog.DebugLevel
	case "info":
		zeroLogLevel = zerolog.InfoLevel
	case "warn":
		zeroLogLevel = zerolog.WarnLevel
	case "fatal":
		zeroLogLevel = zerolog.FatalLevel
	case "panic":
		zeroLogLevel = zerolog.PanicLevel
	case "no":
		zeroLogLevel = zerolog.NoLevel
	case "disabled":
		zeroLogLevel = zerolog.Disabled
	}

	zerolog.SetGlobalLevel(zeroLogLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	utils.Sigs = make(chan os.Signal)
	d := downloader.GetDownloaderInstance(*conf.Global.DownloaderInstanceNum)

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

		if keepRunning {
			return
		}

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
