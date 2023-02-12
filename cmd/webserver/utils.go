package main

import (
	"log"
	"os"

	"github.com/romeq/usva/internal/utils"
)

func setLogWriter(file string) *os.File {
	if file == "" {
		return nil
	}

	fhandle, err := os.OpenFile(file, os.O_WRONLY, 0o644)
	utils.Must(err)

	log.SetOutput(fhandle)
	return fhandle
}
