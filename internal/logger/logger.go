package logger

import (
	"log"
	"time"
)

// TimeTemplate for errors
const timeTemplate = time.RFC822 // "2006.01.02 15:04:05"

func IsError(err error, msg string) bool {
	if err != nil {
		log.Printf("[%s] Error! %s: %s", time.Now().Format(timeTemplate), msg, err)
		return true
	}
	return false
}
