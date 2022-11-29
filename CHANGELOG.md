
<a name="v0.1.0"></a>
## v0.1.0 (2022-11-27)

### Chore

* add codecov

### Docs

* update README.md

### Feat

* add Dockerfile
* add timeout mechanism when download
* skip if file is already downloaded
* add cron job feature
* change config layout
* add download_dir and downloader_instance_num configs
* add command line flags override mechanism
* add config reload mechanism
* fix downloader counter race issue
* update readme
* add config.yaml
* fix test
* change project layout
* change config format to yaml
* replace string key with const
* add downloader counter
* add zerolog

### Fix

* reload config not set log level & unit test
* unit test failure
* main func test fail
* coverage ci
* sync.map range return false issue
* response status 429 Too Many Requests

### Refactor

* wrap errors
* remvoe return err in NewFunc
* change function names
* add interface, instead of struct

### Update

* change project name
* .gitignore
* Makefile

