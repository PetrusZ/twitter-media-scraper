package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
)

type FileDownloader struct {
	FileUrls chan string
	wait sync.WaitGroup
}

func (downloader FileDownloader) Start(workerCount int) {
	for i:=0; i<workerCount; i++ {
		go func() {
			workerId := i
			for fileUrl := range downloader.FileUrls {
				downloader.wait.Add(1)
				fmt.Printf("worker %v downloading file %v\n", workerId, fileUrl)
				err := DownloadFile(fileUrl)
				if err != nil {
					fmt.Println("error downloading: ", err, ": ", fileUrl)
				}
				downloader.wait.Done()
			}
		}()
	}
}

func (downloader FileDownloader) Done() {
	downloader.wait.Wait()
}
func DownloadFile(fileUrl string) error {
	resp, err := http.Get(fileUrl)
	if err != nil  {
		return err
	}
	defer resp.Body.Close()

	parsedUrl, err :=  url.Parse(fileUrl)
	if err != nil {
		return err
	}

	f, err := os.Create("out/"+path.Base(parsedUrl.Path))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
