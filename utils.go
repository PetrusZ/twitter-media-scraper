package main

import "os"

func mkdir(dir string) error {
    _, err := os.Stat(dir)

    if err != nil {
        err := os.Mkdir(dir, os.ModePerm)

        if err != nil {
            return err
        }
    }
    return nil
}
