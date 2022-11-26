package config

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/PetrusZ/twitter-media-scraper/internal/utils"
)

func TestLoad(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	var tests = []struct {
		path     string
		expected bool
	}{
		{basepath + "/configs", false},
		{basepath + "/../../configs", true},
	}

	for _, tt := range tests {
		_, actual := Load(tt.path)
		if !(tt.expected == true && actual == nil) && !(tt.expected == false && actual != nil) {
			t.Errorf("Load(%s): err = %s, expected %s", tt.path, actual, utils.ConvertBoolToString(tt.expected))
		}
	}
}
