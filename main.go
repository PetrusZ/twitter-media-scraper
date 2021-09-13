package main

import (
	"context"
	"flag"
	"fmt"

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

    d := &downloader{info: make(chan tweetInfo)}
    d.Start(16)

    for tweet := range  tweets {
        if tweet.Error != nil {
            return tweet.Error
        }

        url := "https://twitter.com/" + user + "/status/" + tweet.ID

        if *getVideos {
            if tweet.Videos != nil {
                tweetInfo := tweetInfo{user, "", url, Video}
                d.info <- tweetInfo
            }
        }

        if *getPhotos {
            if tweet.Videos == nil {
                for id, url := range tweet.Photos {
                    date := fmt.Sprintf("%d%02d%02d", tweet.TimeParsed.Year(), tweet.TimeParsed.Month(), tweet.TimeParsed.Day())
                    tweetInfo := tweetInfo{user, date + " - " + tweet.ID + "-" + fmt.Sprint(id), url, Photo}
                    d.info <- tweetInfo
                }
            }
        }
    }

    close(d.info)
    d.Wait()

    return err
}
