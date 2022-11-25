package main

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/PetrusZ/twitter-media-scraper/internal/downloader"
	"github.com/PetrusZ/twitter-media-scraper/internal/utils"
)

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

	d := downloader.GetDownloaderInstance(16)
	for _, tt := range tests {
		actual := getUserTweets(tt.user, tt.amount, true, true, d)
		if !(tt.expected == true && actual == nil) && !(tt.expected == false && actual != nil) {
			t.Errorf("getUserTweets(%s, %d): err = %s, expected %s", tt.user, tt.amount, actual, utils.ConvertBoolToString(tt.expected))
		}
	}
}

func TestFlags(T *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	// We manipuate the Args to set them up for the testcases
	// after this test we restore the initial args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	cases := []struct {
		Name string
		Args []string
	}{
		{"flags set", []string{"-configPath", basepath + "/../configs"}},
	}
	for _, tc := range cases {
		// this call is required because otherwise flags panics, if args are set between flag.Parse calls
		flag.CommandLine = flag.NewFlagSet(tc.Name, flag.ExitOnError)
		// we need a value to set Args[0] to, cause flag begins parsing at Args[1]
		os.Args = append([]string{tc.Name}, tc.Args...)
		main()
	}
}
