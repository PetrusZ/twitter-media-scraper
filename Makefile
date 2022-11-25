.PHONY: run build clean

SOURCE = cmd/main.go

TEMP = out/BBCWorld out/wbpictures cmd/out internal/downloader/test internal/utils/test test

PROJECT = twitter-media-scraper

all: run

run:
	go run $(SOURCE)

build:
	go build -o $(PROJECT) $(SOURCE)

test: clean
	go test -v -cover ./...

clean:
	rm -rf $(TEMP)

cleanall:
	rm -rf $(TEMP) $(PROJECT)
