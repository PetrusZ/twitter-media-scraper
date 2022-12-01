package utils

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var testDir = "test"

func setupMkdirAll() {
	MkdirAllFunc = func(string, fs.FileMode) error {
		return errors.New("Can't mkdirAll")
	}
}

func cleanupMkdirAll() {
	MkdirAllFunc = os.MkdirAll
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

	Mkdir(testDir)
	for _, tt := range tests {
		actual := Mkdir(testDir + "/" + tt.dir)
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

	Mkdir(testDir)
	os.RemoveAll("testDir" + "/" + "1")
	for _, tt := range tests {
		actual := MkdirAll(testDir + "/" + tt.dir)
		if !(tt.expected == true && actual == nil) && !(tt.expected == false && actual != nil) {
			t.Errorf("mkdir(%s): err = %s, expected %s", tt.dir, actual, convertBoolToString(tt.expected))
		}
	}

	setupMkdirAll()

	err := MkdirAll("abc")
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

func TestConvertBoolToString(t *testing.T) {
	var tests = []struct {
		b        bool
		expected string
	}{
		{true, "successed"},
		{false, "failed"},
	}

	for _, tt := range tests {
		actual := ConvertBoolToString(tt.b)
		if tt.expected != actual {
			t.Errorf("ConvertBoolToString(%v): actual = %s, expected %s", tt.b, actual, tt.expected)
		}
	}
}

func TestCreate(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected bool
	}{
		{"OK", "foo", true},
		{"Err", "foo", false},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			err := Create(tc.path)
			if tc.expected == true {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}

	err := os.Remove(testCases[0].path)
	require.NoError(t, err)
}
