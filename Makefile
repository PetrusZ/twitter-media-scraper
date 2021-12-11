.PHONY: run build clean

SOURCE = main.go downloader.go utils.go config.go

TEMP = out/BBCWorld out/wbpictures test

PROJECT = twitter-media-scraper

all: run

run:
	go run $(SOURCE)

build:
	go build -o $(PROJECT) $(SOURCE)

test:
	go test

clean:
	rm -rf $(TEMP)

cleanall:
	rm -rf $(TEMP) $(PROJECT)
