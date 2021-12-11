package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
)

var downloaderInstance Downloader
var once sync.Once

type tweetType int

const (
	Photo tweetType = iota
	Video
)

type tweetInfo struct {
	dir  string
	name string
	url  string
}

type Downloader interface {
	Start(int)
	Wait()
	GetInfo() chan tweetInfo
	downloadFile(string, string, string) error
}

type downloader struct {
	info chan tweetInfo
	wg   sync.WaitGroup
}

func NewDownloader() (Downloader, error) {
	return &downloader{info: make(chan tweetInfo)}, nil
}

func GetDownloaderInstance(count int) Downloader {
	once.Do(func() {
		downloaderInstance, _ = NewDownloader()
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
			// log.Printf("workerId %d start\n", workerID)

			defer d.wg.Done()

			for info := range d.info {

				// log.Printf("workerId %d got tweetInfo: dir %s, name %s, url %s\n", workerID, info.dir, info.name, info.url)

				err := d.downloadFile("out/"+info.dir, info.name, info.url)

				if err != nil {
					log.Printf("workerId %d got tweetInfo: dir %s, name %s, url %s\n", workerID, info.dir, info.name, info.url)
					log.Printf("Error: %s", err)
				}

			}
			// log.Printf("workerId %d end\n", workerID)
		})
	}
}

func (d *downloader) Wait() {
	d.wg.Wait()
}

func (d *downloader) downloadFile(dir, name, downloadURL string) error {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	parsedURL, err := url.Parse(downloadURL)
	if err != nil {
		return err
	}

	err = mkdirAll(dir + "/")
	if err != nil {
		return err
	}

	f, err := os.Create(dir + "/" + name + path.Ext(parsedURL.Path))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

/*
func (d *downloader) downloadVideo(dir, url string) error {

	arg := dir + "/%(upload_date)s - %(id)s.%(ext)s"

	cmd := exec.Command("youtube-dl", "-o", arg, url)
	err := cmd.Run()
	if err != nil {
		fmt.Println("cmd error: ", err.Error())
		fmt.Println(cmd.String())
		return err
	}

	return nil
}
*/
