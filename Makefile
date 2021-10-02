.PHONY: run build clean

SOURCE = main.go downloader.go utils.go config.go

TEMP = out/BBCWorld out/wbpictures test

all: run

run:
	go run $(SOURCE)

build:
	go build -o twitter_scraper $(SOURCE)

test:
	go test

clean:
	rm -rf $(TEMP)

cleanall:
	rm -rf $(TEMP) twitter_scraper
