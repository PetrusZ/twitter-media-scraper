package utils

import (
	"fmt"
	"os"
	"runtime/debug"
)

var (
	MkdirAllFunc = os.MkdirAll
	Sigs         chan os.Signal
)

func ConvertBoolToString(b bool) string {
	if b {
		return "successed"
	}
	return "failed"
}

func Mkdir(dir string) error {
	_, err := os.Stat(dir)

	if err != nil {
		err := os.Mkdir(dir, os.ModePerm)

		if err != nil {
			return err
		}
	}
	return nil
}

func MkdirAll(dir string) error {
	_, err := os.Stat(dir)

	if err != nil {
		err := MkdirAllFunc(dir, os.ModePerm)

		if err != nil {
			return err
		}
	}
	return nil
}

func Go(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				debug.PrintStack()
			}
		}()

		fn()
	}()
}
