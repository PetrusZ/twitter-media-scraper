package config

import (
	"github.com/spf13/viper"
)

type configFile struct {
	fileName string
	Configs  []Config
}

type Config struct {
	Global *GlobalConfig `mapstructure:"global"`
	Users  []*UserConfig `mapstructure:"users"`
}

type GlobalConfig struct {
	LogLevel    *string `mapstructure:"log_level"`
	TweetAmount *int    `mapstructure:"tweet_amount"`
	GetVideos   *bool   `mapstructure:"get_videos"`
	GetPhotos   *bool   `mapstructure:"get_photos"`
}

type UserConfig struct {
	UserName    *string `mapstructure:"username"`
	TweetAmount *int    `mapstructure:"tweet_amount"`
	GetVideos   *bool   `mapstructure:"get_videos"`
	GetPhotos   *bool   `mapstructure:"get_photos"`
}

func Load(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

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

	return
}
