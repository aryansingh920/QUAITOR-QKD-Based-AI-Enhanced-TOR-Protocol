package utils

import (
	"log"
	"os"
)

func NewLogger(filename string) *log.Logger {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    return log.New(file, "", log.LstdFlags)
}
