package utils

import (
	"log"
)

// Check logs error and exits program
func Check(err error) {
	if err != nil {
		log.Panic(err.Error())
	}
}
