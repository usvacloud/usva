package main

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/romeq/usva/dbengine"
)

func removeOldFilesWorker(timeUntilRemove time.Duration, uploadDirectory string, cleantrashes bool) {
	for {
		<-time.After(time.Second)

		files, err := dbengine.LastSeenAll()
		if err != nil {
			log.Println("file cleanup process failed", err)
		}
		if cleantrashes {
			go removeJunkWorker(files, uploadDirectory)
		}

		for _, file := range files {
			if time.Now().Before(file.LastSeen.Add(timeUntilRemove)) {
				continue
			}

			err := dbengine.DeleteFile(file.Filename)
			if err != nil {
				log.Println("removeOldFilesWorker:", err)
			}

			err = os.RemoveAll(path.Join(uploadDirectory, file.Filename))
			if err != nil {
				log.Println("removeOldFilesWorker:", err)
			}
		}
	}
}

func removeJunkWorker(files []dbengine.FileLastSeen, uploadDirectory string) {
	fsFiles, err := os.ReadDir(uploadDirectory)
	if err != nil {
		log.Println("removeOldFilesWorker:", err)
	}
	found := false
	for _, direntry := range fsFiles {
		for _, file := range files {
			found = file.Filename == direntry.Name()
			if found {
				break
			}
		}
		if !found {
			err = os.RemoveAll(path.Join(uploadDirectory, direntry.Name()))
			if err != nil {
				log.Println("removeOldFilesWorker:", err)
			}
		}
	}
}
