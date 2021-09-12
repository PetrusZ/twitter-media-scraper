package main

import "testing"

var test_dir = "test"

func TestMkdir(t *testing.T) {
    var tests = []struct {
        dir string
        expected bool
    }{
        {"abc", true},
        {"//", true},
        {"  ", true},
        {"", true},
        {"a/b", false},
        {"a/b/c", false},
    }

    mkdir(test_dir)
    for _, tt := range tests {
        actual := mkdir(test_dir + "/" + tt.dir)
        if  !(tt.expected == true  && actual == nil) && !(tt.expected == false  && actual != nil) {
            t.Errorf("mkdir(%s): err = %s, expected %s", tt.dir, actual, convertBoolToString(tt.expected))
        }
    }
}

func TestDownloadFile(t *testing.T) {
    var tests = []struct {
        dir string
        fileUrl string
        name string
        expected bool
    }{
        {"", "https://www.baidu.com", "https_index", true},
        {"baidu", "", "https_index", false},
        {"baidu", "https://www.baidu.com", "", false},
        {"baidu", "https://www.baidu.com", "https_index", true},
        {"baidu", "http://www.baidu.com", "http_index", true},
        {"baidu", "http://www.baidu.co", "index", false},
    }

    for _, tt := range tests {
        actual := downloadFile(test_dir + "/" + tt.dir, tt.fileUrl, tt.name)
        if  !(tt.expected == true  && actual == nil) && !(tt.expected == false  && actual != nil) {
            t.Errorf("downloadFile(%s, %s, %s): err = %s, expected %s", tt.dir, tt.fileUrl, tt.name, actual, convertBoolToString(tt.expected))
        }
    }
}

func  TestGetUserTweets(t *testing.T) {
    var tests = []struct {
        user string
        amount int
        expected bool
    }{
        {"128j122js,.xzdmcvwe", 2, false},
        {"BBCWorld", 20, true},
        {"BBCWorld", 0, true},
        {"", 2, false},
    }

    for _, tt := range tests {
        actual := getUserTweets(tt.user, tt.amount)
        if  !(tt.expected == true  && actual == nil) && !(tt.expected == false  && actual != nil) {
            t.Errorf("getUserTweets(%s, %d: err = %s, expected %s", tt.user, tt.amount, actual, convertBoolToString(tt.expected))
        }
    }
}

func convertBoolToString(b bool) string{
    if b {
        return "successed"
    } else {
        return "failed"
    }
}
