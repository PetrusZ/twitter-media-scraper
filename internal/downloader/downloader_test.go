package downloader

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/PetrusZ/twitter-media-scraper/internal/config"
	"github.com/PetrusZ/twitter-media-scraper/internal/utils"
	"github.com/stretchr/testify/require"
)

var testDir = "test"

func setupMkdirAll() {
	utils.MkdirAllFunc = func(string, fs.FileMode) error {
		return errors.New("Can't mkdirAll")
	}
}

func cleanupMkdirAll() {
	utils.MkdirAllFunc = os.MkdirAll
}

func TestDownloadFile(t *testing.T) {
	var tests = []struct {
		dir      string
		fileURL  string
		name     string
		expected bool
	}{
		{"", "https://www.baidu.com", "https_index", true},
		{"baidu", "", "https_index", false},
		{"baidu", "https://www.baidu.com", "", false},
		{"baidu", "https://www.baidu.com", "https_index", true},
		{"baidu", "http://www.baidu.com", "http_index", true},
		{"baidu", "http://www.baidu.co", "index", false},
		{"baidu", "baidu", "index", false},
	}

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	configPath := basepath + "/../../configs"
	_, err := config.Load(configPath)
	require.NoError(t, err)

	d := GetDownloaderInstance(16)
	for _, tt := range tests {
		actual := d.downloadFile(testDir+"/"+tt.dir, tt.name, tt.fileURL)
		if !(tt.expected == true && actual == nil) && !(tt.expected == false && actual != nil) {
			t.Errorf("downloadFile(%s, %s, %s): err = %v, expected %s", tt.dir, tt.fileURL, tt.name, actual, utils.ConvertBoolToString(tt.expected))
		}
	}

	setupMkdirAll()

	err = d.downloadFile(testDir+"/testDownloadFile", "index", "http://www.baidu.com")
	require.Error(t, err)

	cleanupMkdirAll()

	os.RemoveAll(testDir)
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
			t.Errorf("downloadFile(%s, %s, %s): err = %s, expected %s", tt.dir, tt.fileURL, tt.name, actual, utils.ConvertBoolToString(tt.expected))
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
