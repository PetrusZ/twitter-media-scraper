.PHONY: run build clean

SOURCE = cmd/main.go

TEMP = out/BBCWorld out/wbpictures internal/downloader/test internal/utils/test test

PROJECT = twitter-media-scraper

all: run

run:
	go run $(SOURCE) --keep_running=false

build:
	go build -o $(PROJECT) $(SOURCE)

build_all:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(PROJECT).linux-amd64 $(SOURCE)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o $(PROJECT).linux-arm $(SOURCE)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(PROJECT).linux-arm64 $(SOURCE)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(PROJECT).windows-amd64 $(SOURCE)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(PROJECT).darwin-amd64 $(SOURCE)

test: clean
	go test -v -cover ./...

clean:
	rm -rf $(TEMP) $(PROJECT)*

clean_debug:
	rm -rf cmd/out

cleanall:
	rm -rf $(TEMP) $(PROJECT)

docker_mac:
	docker buildx build --push -f build/package/Dockerfile --platform linux/amd64,linux/arm64 -t patrickz07/$(PROJECT):latest .