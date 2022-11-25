package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"

	"github.com/rs/zerolog/log"
)

var downloaderInstance Downloader
var once sync.Once

type tweetType int

const (
	TweetTypePhoto tweetType = iota
	TweetTypeVideo

	CounterKeyVideo = "video"
	CounterKeyPhoto = "photo"
	CounterKeyTotal = "total"
)

type tweetInfo struct {
	user      string
	dir       string
	name      string
	url       string
	tweetType tweetType
}

type Downloader interface {
	Start(int)
	Wait()
	GetInfo() chan tweetInfo
	PrintCounter()
	downloadFile(string, string, string) error
}

type downloader struct {
	info    chan tweetInfo
	wg      sync.WaitGroup
	counter sync.Map
}

func NewDownloader() Downloader {
	return &downloader{info: make(chan tweetInfo)}
}

func GetDownloaderInstance(count int) Downloader {
	once.Do(func() {
		downloaderInstance = NewDownloader()
		downloaderInstance.Start(count)
	})
	return downloaderInstance
}

func (d *downloader) GetInfo() chan tweetInfo {
	return d.info
}

func (d *downloader) Start(count int) {
	for i := 0; i < count; i++ {
		workerID := i
		d.wg.Add(1)
		Go(func() {
			log.Debug().Msgf("workerId %d start\n", workerID)

			defer d.wg.Done()

			for info := range d.info {

				log.Debug().Msgf("workerId %d got tweetInfo: dir %s, name %s, url %s\n", workerID, info.dir, info.name, info.url)

				err := d.downloadFile("out/"+info.dir, info.name, info.url)
				if err != nil {
					log.Error().Msgf("workerId %d got tweetInfo: dir %s, name %s, url %s\n", workerID, info.dir, info.name, info.url)
					log.Error().Msgf("Error: %s", err)
				} else {
					d.increaseCounter(info.user, info.tweetType)
				}

			}
			log.Debug().Msgf("workerId %d end\n", workerID)
		})
	}
}

func (d *downloader) Wait() {
	d.wg.Wait()
}

func (d *downloader) downloadFile(dir, name, downloadURL string) error {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("http.Get(%s) error: %w", downloadURL, err)
	}
	defer resp.Body.Close()

	parsedURL, err := url.Parse(downloadURL)
	if err != nil {
		return fmt.Errorf("url.Parse(%s) error: %w", downloadURL, err)
	}

	err = mkdirAll(dir + "/")
	if err != nil {
		return fmt.Errorf("mkdirAll(%s) error: %w", dir+"/", err)
	}

	f, err := os.Create(dir + "/" + name + path.Ext(parsedURL.Path))
	if err != nil {
		return fmt.Errorf("os.Create(%s) error: %w", dir+"/"+name+path.Ext(parsedURL.Path), err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("io.Copy error: %w", err)
	}

	return nil
}

func (d *downloader) increaseCounter(user string, tweetType tweetType) {
	count := 1
	totalCount := 1
	counter := map[string]int{
		CounterKeyPhoto: 0,
		CounterKeyVideo: 0,
		CounterKeyTotal: 0,
	}

	subKey := CounterKeyTotal
	if tweetType == TweetTypeVideo {
		subKey = CounterKeyVideo
	} else if tweetType == TweetTypePhoto {
		subKey = CounterKeyPhoto
	}

	counterIntf, ok := d.counter.Load(user)
	if ok {
		// can get current count
		counter, ok = counterIntf.(map[string]int)
		if !ok {
			log.Error().Msgf("counterIntf = %v convert to map[string]int failed", counter)
			return
		}

		countInt, ok := counter[subKey]
		if !ok {
			countInt = 0
		}
		count = countInt + 1

		totalCountInt, ok := counter[CounterKeyTotal]
		if !ok {
			totalCountInt = 0
		}
		totalCount = totalCountInt + 1
	}

	counter[subKey] = count
	counter[CounterKeyTotal] = totalCount
	d.counter.Store(user, counter)
}

func (d *downloader) PrintCounter() {
	d.counter.Range(func(user, subKey interface{}) bool {
		userStr, ok := user.(string)
		if !ok {
			log.Error().Msgf("user %v conver to string failed", user)
			return true
		}

		counter, ok := subKey.(map[string]int)
		if !ok {
			log.Error().Msgf("user %v conver to map[string]int failed", subKey)
			return true
		}

		log.Info().Msgf("user %s downloaded %d photo(s), %d video(s), %d total",
			userStr, counter[CounterKeyPhoto], counter[CounterKeyVideo], counter[CounterKeyTotal])

		return true
	})
}
