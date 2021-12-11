package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type BodyReader func(io.Reader) ([]byte, error)
type RespUnmarshaller func([]byte, interface{}) error

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

var readAllFunc = ioutil.ReadAll
var unMarshalFunc = json.Unmarshal

func NewConfigFile(fileName string) ConfigFile {
	return &configFile{
		fileName: fileName,
	}
}

func (c *configFile) GetConfigs() []Config {
	return c.Configs
}

func (c *configFile) Load() error {
	jsonFile, err := os.Open(c.fileName)
	if err != nil {
		return fmt.Errorf("os.Open(%s) error: %w", c.fileName, err)
	}
	defer jsonFile.Close()

	byteValue, err := readAllFunc(jsonFile)
	if err != nil {
		return fmt.Errorf("readAllFunc error: %w", err)
	}

	err = unMarshalFunc(byteValue, &c.Configs)
	if err != nil {
		return fmt.Errorf("unMarshalFunc error: %w", err)
	}

	return nil
}
