package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"sync"
)

type tweetType int

const (
    Photo tweetType = iota
    Video
)

type tweetInfo struct {
    dir string
    name string
    url string
    tweetType tweetType
}

type downloader struct {
    info chan tweetInfo
    wg sync.WaitGroup
}

func (d *downloader) Start (count int) {
    for i := 0; i < count; i++ {
        workerId := i
        go func() {
            // log.Printf("workerId %d start\n", workerId)

            for info := range d.info {
                d.wg.Add(1)

                // log.Printf("workerId %d got tweetInfo: dir %s, name %s, url %s\n", workerId, info.dir, info.name, info.url)

                var err error
                if info.tweetType == Video {
                    err = d.downloadVideo(info.dir, info.url)
                } else if info.tweetType == Photo {
                    err = d.downloadFile(info.dir, info.name, info.url + "?format=jpg&name=orig")
                }

                if err != nil {
                    log.Printf("workerId %d got tweetInfo: dir %s, name %s, url %s\n", workerId, info.dir, info.name, info.url)
                    log.Printf("Error: %s", err)
                }

                d.wg.Done()
            }
            // log.Printf("workerId %d end\n", workerId)
        }()
    }
}

func (d *downloader) Wait(){
    d.wg.Wait()
}

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

func (d *downloader) downloadFile(dir, name, downloadUrl string) error {
    resp, err := http.Get(downloadUrl)
    if err != nil  {
        return err
    }
    defer resp.Body.Close()

    parsedUrl, err :=  url.Parse(downloadUrl)
    if err != nil {
        return err
    }

    err = mkdir(dir)
    if err != nil {
        return err
    }

    f, err := os.Create(dir + "/" + name + path.Ext(parsedUrl.Path))
    if err != nil {
        return err
    }
    defer f.Close()

    _, err = io.Copy(f, resp.Body)
    return err
}
