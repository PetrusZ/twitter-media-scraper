package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type ConfigFile interface {
	Load() error
	GetConfigs() []Config
}

type configFile struct {
	fileName string
	Configs  []Config
}

type Config struct {
	UserName    string `json:"user_name"`
	TweetAmount int    `json:"tweet_amount"`
	GetVideos   bool   `json:"get_videos"`
	GetPhotos   bool   `json:"get_photos"`
}

func NewConfigFile(fileName string) (ConfigFile, error) {
	return &configFile{fileName: fileName}, nil
}

func (c *configFile) GetConfigs() []Config {
	return c.Configs
}

func (c *configFile) Load() error {
	jsonFile, err := os.Open(c.fileName)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, &c.Configs)
	if err != nil {
		return err
	}

	return nil
}
