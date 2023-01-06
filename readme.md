# twitter-media-scraper

[![codecov](https://codecov.io/gh/PetrusZ/twitter-media-scraper/branch/main/graph/badge.svg)](https://codecov.io/gh/PetrusZ/twitter-media-scraper)

# Intro

Scrape/Craw twitter users' pictures and videos by username.

# Feature

* Automaticly reload config when config file change
* Support cron job to schedule download
* Support deploy by docker
* Skip downloading if file already exist

# Usage

## Build from source

1. mv `configs/config.example.yaml` to `configs/config.yaml`
1. Edit `config.yaml`
2. `make run`

## Docker

``` sh
docker run -d \
  --name twitter-media-scraper \
  -v /etc/timezone:/etc/timezone:ro \
  -v /etc/localtime:/etc/localtime:ro \
  -v /path/to/out:/cmd/out \
  -v /path/to/configs:/cmd/configs \
  --restart=always \
  patrickz07/twitter-media-scraper:latest
```