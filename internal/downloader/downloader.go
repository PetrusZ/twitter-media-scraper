package downloader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"

	"github.com/PetrusZ/twitter-media-scraper/internal/utils"
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

type TweetInfo struct {
	User      string
	Dir       string
	Name      string
	URL       string
	TweetType tweetType
}

type Downloader interface {
	Start(int)
	Init()
	Wait()
	GetInfo() chan TweetInfo
	PrintCounter()
	downloadFile(string, string, string) error
}

type downloader struct {
	info    chan TweetInfo
	wg      sync.WaitGroup
	counter sync.Map
}

func NewDownloader() Downloader {
	return &downloader{info: make(chan TweetInfo)}
}

func GetDownloaderInstance(count int) Downloader {
	once.Do(func() {
		downloaderInstance = NewDownloader()
		downloaderInstance.Start(count)
	})
	return downloaderInstance
}

func (d *downloader) GetInfo() chan TweetInfo {
	return d.info
}

func (d *downloader) Start(count int) {
	for i := 0; i < count; i++ {
		workerID := i
		d.wg.Add(1)
		utils.Go(func() {
			log.Debug().Msgf("workerId %d start", workerID)

			defer d.wg.Done()

			for info := range d.info {

				log.Debug().Msgf("workerId %d got tweetInfo: dir %s, name %s, url %s", workerID, info.Dir, info.Name, info.URL)

				err := d.downloadFile("out/"+info.Dir, info.Name, info.URL)
				if err != nil {
					log.Error().Msgf("workerId %d got tweetInfo: dir %s, name %s, url %s", workerID, info.Dir, info.Name, info.URL)
					log.Error().Msgf("Error: %s", err)
				} else {
					d.increaseCounter(info.User, info.TweetType)
				}

			}
			log.Debug().Msgf("workerId %d end", workerID)
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

	err = utils.MkdirAll(dir + "/")
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

	var counter sync.Map
	counter.Store(CounterKeyPhoto, 0)
	counter.Store(CounterKeyVideo, 0)
	counter.Store(CounterKeyTotal, 0)

	subKey := CounterKeyTotal
	if tweetType == TweetTypeVideo {
		subKey = CounterKeyVideo
	} else if tweetType == TweetTypePhoto {
		subKey = CounterKeyPhoto
	}

	counterIntf, ok := d.counter.Load(user)
	if ok {
		// can get current count
		counter, ok = counterIntf.(sync.Map)
		if !ok {
			log.Error().Msgf("counterIntf = %v convert to sync.Map failed", counter)
			return
		}

		countInt := 0
		countIntf, ok := counter.Load(subKey)
		if ok {
			countInt, ok = countIntf.(int)
			if !ok {
				log.Error().Msgf("countIntf = %v convert to int failed", countIntf)
				return
			}
		}
		count = countInt + 1

		totalCountInt := 0
		totalCountIntf, ok := counter.Load(CounterKeyTotal)
		if ok {
			totalCountInt, ok = totalCountIntf.(int)
			if !ok {
				log.Error().Msgf("countIntf = %v convert to int failed", totalCountIntf)
				return
			}
		}
		totalCount = totalCountInt + 1
	}

	counter.Store(subKey, count)
	counter.Store(CounterKeyTotal, totalCount)
	d.counter.Store(user, counter)
}

func (d *downloader) PrintCounter() {
	d.counter.Range(func(user, subKey interface{}) bool {
		userStr, ok := user.(string)
		if !ok {
			log.Error().Msgf("user %v conver to string failed", user)
			return true
		}

		counter, ok := subKey.(sync.Map)
		if !ok {
			log.Error().Msgf("user %v conver to sync.Map failed", subKey)
			return true
		}

		photoCount, ok := counter.Load(CounterKeyPhoto)
		if !ok {
			log.Error().Msg("counter.Load(CounterKeyPhoto) failed")
		}

		videoCount, ok := counter.Load(CounterKeyVideo)
		if !ok {
			log.Error().Msg("counter.Load(CounterKeyVideo) failed")
		}

		totalCount, ok := counter.Load(CounterKeyTotal)
		if !ok {
			log.Error().Msg("counter.Load(CounterKeytotal) failed")
		}

		log.Info().Msgf("user %s downloaded %d photo(s), %d video(s), %d total",
			userStr, photoCount, videoCount, totalCount)

		return true
	})
}

func (d *downloader) Init() {
	d.info = make(chan TweetInfo)
	d.counter.Range(func(user, subKey interface{}) bool {
		d.counter.Delete(user)
		return true
	})
}
