package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

type BodyReader func(io.Reader) ([]byte, error)
type RespUnmarshaller func([]byte, interface{}) error

type ConfigFile interface {
	Load() error
	GetConfigs() []Config
	SetBodyReader(BodyReader)
	SetUnmarshaller(RespUnmarshaller)
}

type configFile struct {
	fileName     string
	Configs      []Config
	bodyReader   BodyReader
	unMarshaller RespUnmarshaller
}

type Config struct {
	UserName    string `json:"user_name"`
	TweetAmount int    `json:"tweet_amount"`
	GetVideos   bool   `json:"get_videos"`
	GetPhotos   bool   `json:"get_photos"`
}

func NewConfigFile(fileName string) (ConfigFile, error) {
	return &configFile{
		fileName:     fileName,
		bodyReader:   ioutil.ReadAll,
		unMarshaller: json.Unmarshal,
	}, nil
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

	byteValue, err := c.bodyReader(jsonFile)
	if err != nil {
		return err
	}

	err = c.unMarshaller(byteValue, &c.Configs)
	if err != nil {
		return err
	}

	return nil
}

func (c *configFile) SetBodyReader(bodyReader BodyReader) {
	c.bodyReader = bodyReader
}

func (c *configFile) SetUnmarshaller(unMarshaller RespUnmarshaller) {
	c.unMarshaller = unMarshaller
}
