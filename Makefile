.PHONY: run build clean

run:
	go run ./main.go ./downloader.go ./utils.go ./config.go

build:
	go build -o twitter_scraper ./main.go ./downloader.go ./utils.go ./config.go

test:
	go test

clean:
	rm -rf BBCWorld test

clean_all:
	rm -rf BBCWorld test twitter_scraper
