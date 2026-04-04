package main

import (
	"log"
	"os"

	"github.com/usvacloud/usva/internal/utils"
)

func setLogWriter(file string) *os.File {
	if file == "" {
		return nil
	}

	handle, err := os.OpenFile(file, os.O_WRONLY, 0o644)
	utils.MustBeNil(err)

	log.SetOutput(handle)
	return handle
}
