package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

var  (
    twitterUser = flag.String("user", "", "the twitter user to get files from")
    tweetAmount = flag.Int("amount", 1000, "amount of tweets to get content from")
    getVideos = flag.Bool("videos", true, "download videos from tweets")
    getPhotos = flag.Bool("photos", true, "download photos from tweets")
)

func main() {
    if *twitterUser != "" {
        getUserTweets(*twitterUser, *tweetAmount)
    } else {
        fmt.Println("No twitter user")
    }
}

func getUserTweets(user string, amount int) (err error) {
    scraper := twitterscraper.New()


    tweets := scraper.GetTweets(context.Background(), user, amount)

    if tweets == nil {
        err = mkdir(user)
        if err != nil {
            return err
        }
    }

    for tweet := range  tweets {
        if tweet.Error != nil {
            return tweet.Error
        }

        url := "https://twitter.com/" + user + "/status/" + tweet.ID

        if *getVideos {
            if tweet.Videos != nil {
                arg := user + "/%(upload_date)s - %(id)s.%(ext)s"

                cmd := exec.Command("youtube-dl", "-o", arg, url)
                err := cmd.Run()
                if err != nil {
                    fmt.Println("cmd error: ", err.Error())
                    fmt.Println(cmd.String())
                    return err
                }
            }
        }

        if *getPhotos {
            if tweet.Videos == nil {
                for id, photo := range tweet.Photos {
                    date := fmt.Sprintf("%d%02d%02d", tweet.TimeParsed.Year(), tweet.TimeParsed.Month(), tweet.TimeParsed.Day())
                    downloadFile(user, photo, date + " - " + tweet.ID + "-" + fmt.Sprint(id))
                }
            }
        }
    }

    return err
}

func downloadFile(dir, fileUrl, name string) error {
    resp, err := http.Get(fileUrl)
    if err != nil  {
        return err
    }
    defer resp.Body.Close()

    parsedUrl, err :=  url.Parse(fileUrl)
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

func mkdir(dir string) error {
    _, err := os.Stat(dir)

    if err != nil {
        err := os.Mkdir(dir, os.ModePerm)

        if err != nil {
            return err
        }
    }
    return nil
}
