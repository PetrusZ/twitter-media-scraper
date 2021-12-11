package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"io/fs"
	"os"
	"testing"
)

var testDir = "test"

func setupMkdirAll() {
	mkdirAllFunc = func(string, fs.FileMode) error {
		return errors.New("Can't mkdirAll")
	}
}

func cleanupMkdirAll() {
	mkdirAllFunc = os.MkdirAll
}

func convertBoolToString(b bool) string {
	if b {
		return "successed"
	}
	return "failed"
}

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
	os.RemoveAll("testDir" + "/" + "1")
	for _, tt := range tests {
		actual := mkdirAll(testDir + "/" + tt.dir)
		if !(tt.expected == true && actual == nil) && !(tt.expected == false && actual != nil) {
			t.Errorf("mkdir(%s): err = %s, expected %s", tt.dir, actual, convertBoolToString(tt.expected))
		}
	}

	setupMkdirAll()

	err := mkdirAll("abc")
	if err == nil {
		t.Error("mkdirAll expected has err, but got nil")
	}

	cleanupMkdirAll()
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

	d := GetDownloaderInstance(16)
	for _, tt := range tests {
		actual := d.downloadFile(testDir+"/"+tt.dir, tt.name, tt.fileURL)
		if !(tt.expected == true && actual == nil) && !(tt.expected == false && actual != nil) {
			t.Errorf("downloadFile(%s, %s, %s): err = %s, expected %s", tt.dir, tt.fileURL, tt.name, actual, convertBoolToString(tt.expected))
		}
	}

	setupMkdirAll()

	err := d.downloadFile(testDir+"/testDownloadFile", "index", "http://www.baidu.com")
	if err == nil {
		t.Error("downloadFile expected err, but got nil")
	}

	cleanupMkdirAll()
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
		readAllFunc = func(io.Reader) ([]byte, error) {
			return nil, errors.New("")
		}
		err := config.Load("config.json")
		if err == nil {
			t.Errorf("Load(): err = %s, expected err", err)
		}
	})

	t.Run("error on unMarshaller", func(t *testing.T) {
		readAllFunc = io.ReadAll
		unMarshalFunc = func([]byte, interface{}) error {
			return errors.New("")
		}
		err := config.Load("config.json")
		if err == nil {
			t.Errorf("Load(): err = %s, expected err", err)
		}
	})

	readAllFunc = io.ReadAll
	unMarshalFunc = json.Unmarshal
}

func TestFlags(T *testing.T) {
	// We manipuate the Args to set them up for the testcases
	// after this test we restore the initial args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	cases := []struct {
		Name string
		Args []string
	}{
		// {"flags set", []string{"-configFile", "dontExits"}},
		{"flags not set", []string{""}},
	}
	for _, tc := range cases {
		// this call is required because otherwise flags panics, if args are set between flag.Parse calls
		flag.CommandLine = flag.NewFlagSet(tc.Name, flag.ExitOnError)
		// we need a value to set Args[0] to, cause flag begins parsing at Args[1]
		os.Args = append([]string{tc.Name}, tc.Args...)
		main()
	}
}
