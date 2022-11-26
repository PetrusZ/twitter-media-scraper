package config

import (
	"syscall"

	"github.com/PetrusZ/twitter-media-scraper/internal/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var config Config

type Config struct {
	LogLevel              *string       `mapstructure:"log_level"`
	KeepRunning           *bool         `mapstructure:"keep_running"`
	DownloadDir           *string       `mapstructure:"download_dir"`
	DownloaderInstanceNum *int          `mapstructure:"downloader_instance_num"`
	Cron                  *string       `mapstructure:"cron"`
	Global                *GlobalConfig `mapstructure:"global"`
	Users                 []*UserConfig `mapstructure:"users"`
}

type GlobalConfig struct {
	TweetAmount *int  `mapstructure:"tweet_amount"`
	GetVideos   *bool `mapstructure:"get_videos"`
	GetPhotos   *bool `mapstructure:"get_photos"`
}

type UserConfig struct {
	UserName    *string `mapstructure:"username"`
	TweetAmount *int    `mapstructure:"tweet_amount"`
	GetVideos   *bool   `mapstructure:"get_videos"`
	GetPhotos   *bool   `mapstructure:"get_photos"`
}

func Get() Config {
	return config
}

func Load(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	config = Config{}
	err := viper.ReadInConfig()
	if err != nil {
		log.Error().Err(err).Msg("read in config error")
		return config, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Error().Err(err).Msg("unmarshal config error")
		return config, err
	}

	for _, user := range config.Users {
		if user.TweetAmount == nil {
			user.TweetAmount = config.Global.TweetAmount
		}

		if user.GetPhotos == nil {
			user.GetPhotos = config.Global.GetPhotos
		}

		if user.GetVideos == nil {
			user.GetVideos = config.Global.GetVideos
		}
	}

	return config, nil
}

func Watch() {
	viper.OnConfigChange(func(e fsnotify.Event) {
		utils.Sigs <- syscall.SIGHUP
	})
	viper.WatchConfig()
}
