package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type ConfigFile struct {
    fileName string
    Configs []Config
}

type Config struct {
    UserName string `json:"user_name"`
    TweetAmount int `json:"tweet_amount"`
    GetVideos bool `json:"get_videos"`
    GetPhotos bool `json:"get_photos"`
}

func NewConfigFile() *ConfigFile {
    return &ConfigFile{}
}

func (c *ConfigFile) Load(name string) error {
    jsonFile, err := os.Open(name)
    if err != nil {
        return err
    }
    defer jsonFile.Close()
    c.fileName = name

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
