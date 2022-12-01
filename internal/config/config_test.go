package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/PetrusZ/twitter-media-scraper/internal/utils"
	"github.com/stretchr/testify/require"
)

func createWrongFormatConfig(t *testing.T, name string) {
	err := utils.Create(name)
	require.NoError(t, err)

	data := []byte("foo: bar\nbar: foo")
	err = os.WriteFile(name, data, 0644)
	require.NoError(t, err)
}

func TestLoad(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	wrongFormatConfigName := "foo_config.yaml"
	createWrongFormatConfig(t, wrongFormatConfigName)
	defer os.Remove(wrongFormatConfigName)

	var tests = []struct {
		name     string
		path     string
		expected bool
	}{
		{"Not found", basepath + "/configs", false},
		{"Wrong format", wrongFormatConfigName, false},
		{"OK", basepath + "/../../configs", true},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			_, actual := Load(tt.path)
			if !(tt.expected == true && actual == nil) && !(tt.expected == false && actual != nil) {
				t.Errorf("Load(%s): err = %s, expected %s", tt.path, actual, utils.ConvertBoolToString(tt.expected))
			}
		})
	}
}

func TestGet(t *testing.T) {
	conf := Get()
	require.NotNil(t, conf)
}

func TestWatch(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	conf, err := Load(basepath + "/../../configs")
	require.NoError(t, err)
	require.NotEmpty(t, conf)

	Watch()
}
