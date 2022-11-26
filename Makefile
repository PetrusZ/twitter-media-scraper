.PHONY: run build clean

SOURCE = cmd/main.go

TEMP = out/BBCWorld out/wbpictures internal/downloader/test internal/utils/test test

PROJECT = twitter-media-scraper

all: run

run:
	go run $(SOURCE) --keep_running=false

build:
	go build -o $(PROJECT) $(SOURCE)

test: clean
	go test -v -cover ./...

clean:
	rm -rf $(TEMP)

clean_debug:
	rm -rf cmd/out

cleanall:
	rm -rf $(TEMP) $(PROJECT)

docker_mac:
	docker buildx build --push -f build/package/Dockerfile --platform linux/amd64,linux/arm64 -t patrickz07/$(PROJECT):latest .