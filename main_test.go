package main

import (
	"errors"
	"io"
	"os"
	"testing"
)

var testDir = "test"

func TestMkdir(t *testing.T) {
	var tests = []struct {
		dir      string
		expected bool
	}{
		{"abc", true},
		{"//", true},
		{"  ", true},
		{"", true},
		{"a/b", false},
		{"a/b/c", false},
	}

	mkdir(testDir)
	for _, tt := range tests {
		actual := mkdir(testDir + "/" + tt.dir)
		if !(tt.expected == true && actual == nil) && !(tt.expected == false && actual != nil) {
			t.Errorf("mkdir(%s): err = %s, expected %s", tt.dir, actual, convertBoolToString(tt.expected))
		}
	}
}

func TestMkdirAll(t *testing.T) {
	var tests = []struct {
		dir      string
		expected bool
	}{
		{"123", true},
		{" /f213/", true},
		{"  ", true},
		{"", true},
		{"1/a/b", true},
		{"1/a/b/c", true},
		{"1/a/b/c/d", true},
		{"1/a/b/c/d/e", true},
	}

	mkdir(testDir)
	for _, tt := range tests {
		actual := mkdirAll(testDir + "/" + tt.dir)
		if !(tt.expected == true && actual == nil) && !(tt.expected == false && actual != nil) {
			t.Errorf("mkdir(%s): err = %s, expected %s", tt.dir, actual, convertBoolToString(tt.expected))
		}
	}
}

func TestGo(t *testing.T) {
	f1 := func() {
	}

	f2 := func() {
		panic("panic test")
	}
	var tests = []struct {
		fn       func()
		expected bool
	}{
		{f1, true},
		{f2, true},
	}

	for _, tt := range tests {
		Go(tt.fn)
	}
}

/*
func TestDownloadVideo(t *testing.T) {
	var tests = []struct {
		dir      string
		fileURL  string
		expected bool
	}{
		{"BBCWorld", "https://twitter.com/BBCWorld/status/1439286551733800960", true},
		{"BBCWorld", "https://twitter.com/BBCWorld/status/14392865517338", false},
	}

	d := GetDownloaderInstance()
	for _, tt := range tests {
		actual := d.downloadVideo(testDir+"/"+tt.dir, tt.fileURL)
		if !(tt.expected == true && actual == nil) && !(tt.expected == false && actual != nil) {
			t.Errorf("downloadFile(%s, %s): err = %s, expected %s", tt.dir, tt.fileURL, actual, convertBoolToString(tt.expected))
		}
	}
}
*/

func TestDownloadFile(t *testing.T) {
	var tests = []struct {
		dir      string
		name     string
		fileURL  string
		expected bool
	}{
		{"", "https://www.baidu.com", "https_index", true},
		{"baidu", "", "https_index", false},
		{"baidu", "https://www.baidu.com", "", false},
		{"baidu", "https://www.baidu.com", "https_index", true},
		{"baidu", "http://www.baidu.com", "http_index", true},
		{"baidu", "http://www.baidu.co", "index", false},
		{"baidu", "http://www.baidu", "index", false},
	}

	d := GetDownloaderInstance(16)
	for _, tt := range tests {
		actual := d.downloadFile(testDir+"/"+tt.dir, tt.fileURL, tt.name)
		if !(tt.expected == true && actual == nil) && !(tt.expected == false && actual != nil) {
			t.Errorf("downloadFile(%s, %s, %s): err = %s, expected %s", tt.dir, tt.fileURL, tt.name, actual, convertBoolToString(tt.expected))
		}
	}
}

func TestParallelDownloadFile(t *testing.T) {
	var tests = []struct {
		dir      string
		fileURL  string
		name     string
		expected bool
	}{
		{"Paralle", "https://t1.huishahe.com/uploads/tu/zyf/tt/20160520/erx0a4ooid2.jpg", "bigPic1", true},
		{"Paralle", "https://t1.huishahe.com/uploads/tu/zyf/tt/20160520/erx0a4ooid2.jpg", "bigPic2", true},
		{"Paralle", "https://t1.huishahe.com/uploads/tu/zyf/tt/20160520/erx0a4ooid2.jpg", "bigPic3", true},
		{"Paralle", "https://t1.huishahe.com/uploads/tu/zyf/tt/20160520/erx0a4ooid2.jpg", "bigPic4", true},
		{"Paralle", "https://t1.huishahe.com/uploads/tu/zyf/tt/20160520/erx0a4ooid2.jpg", "bigPic5", true},
		{"Paralle", "https://t1.huishahe.com/uploads/tu/zyf/tt/20160520/erx0a4ooid2.jpg", "bigPic6", true},
		{"Paralle", "https://t1.huishahe.com/uploads/tu/zyf/tt/20160520/erx0a4ooid2.jpg", "bigPic7", true},
		{"Paralle", "https://t1.huishahe.com/uploads/tu/zyf/tt/20160520/erx0a4ooid2.jpg", "bigPic8", true},
	}

	d := GetDownloaderInstance(16)
	var fileSize int64
	for _, tt := range tests {
		actual := d.downloadFile(testDir+"/"+tt.dir, tt.name, tt.fileURL)
		if !(tt.expected == true && actual == nil) && !(tt.expected == false && actual != nil) {
			t.Errorf("downloadFile(%s, %s, %s): err = %s, expected %s", tt.dir, tt.fileURL, tt.name, actual, convertBoolToString(tt.expected))
		}

		file, err := os.Stat(testDir + "/" + tt.dir + "/" + tt.name + ".jpg")
		if err != nil {
			t.Errorf("downloadFile(%s, %s, %s): downloadFile not exist", tt.dir, tt.fileURL, tt.name)
		}

		size := file.Size()
		if fileSize == 0 {
			fileSize = size
		}

		if fileSize != size {
			t.Errorf("downloadFile(%s, %s, %s): downloadFile size not equal, one is %d, another is %d", tt.dir, tt.fileURL, tt.name, fileSize, size)
		}
	}
}

func TestGetUserTweets(t *testing.T) {
	var tests = []struct {
		user     string
		amount   int
		expected bool
	}{
		{"128j122js,.xzdmcvwe", 50, false},
		{"BBCWorld", 50, true},
		{"BBCWorld", 0, true},
		{"wbpictures", 50, true},
		{"", 50, false},
	}

	d := GetDownloaderInstance(16)
	for _, tt := range tests {
		actual := getUserTweets(tt.user, tt.amount, true, true, d)
		if !(tt.expected == true && actual == nil) && !(tt.expected == false && actual != nil) {
			t.Errorf("getUserTweets(%s, %d): err = %s, expected %s", tt.user, tt.amount, actual, convertBoolToString(tt.expected))
		}
	}
}

func TestLoad(t *testing.T) {
	var tests = []struct {
		name     string
		expected bool
	}{
		{"config.json", true},
		{"abc.json", false},
	}

	for _, tt := range tests {
		config, _ := NewConfigFile(tt.name)
		actual := config.Load()
		if !(tt.expected == true && actual == nil) && !(tt.expected == false && actual != nil) {
			t.Errorf("Load(%s): err = %s, expected %s", tt.name, actual, convertBoolToString(tt.expected))
		}
	}
	t.Run("error on bodyReader", func(t *testing.T) {
		bodyReader = func(io.Reader) ([]byte, error) {
			return nil, errors.New("")
		}
		err := config.Load("config.json")
		if err == nil {
			t.Errorf("Load(): err = %s, expected err", err)
		}
	})

	t.Run("error on unMarshaller", func(t *testing.T) {
		bodyReader = io.ReadAll
		unMarshaller = func([]byte, interface{}) error {
			return errors.New("")
		}
		err := config.Load("config.json")
		if err == nil {
			t.Errorf("Load(): err = %s, expected err", err)
		}
	})
}

func convertBoolToString(b bool) string {
	if b {
		return "successed"
	}
	return "failed"
}
