package api

import (
	"context"
	"log"
	"os"
	"path"
	"time"

	"github.com/romeq/usva/internal/db"
)

func (s *Server) RemoveOldFilesWorker(timeUntilRemove time.Duration, uploadDirectory string, cleantrashes bool) {
	log.Println("executed RemoveOldFilesWorker")
	workContext := context.Background()

	files, err := s.db.GetLastSeenAll(workContext)
	if err != nil {
		log.Println("file cleanup process failed", err)
		return
	}

	if cleantrashes {
		go s.RemoveJunkWorker(files, uploadDirectory)
	}

	for _, file := range files {
		if time.Now().Before(file.LastSeen.Add(timeUntilRemove)) {
			continue
		}

		err := s.db.DeleteFile(workContext, file.FileUuid)
		if err != nil {
			log.Println("RemoveOldFilesWorker:", err)
			return
		}

		err = os.RemoveAll(path.Join(uploadDirectory, file.FileUuid))
		if err != nil {
			log.Println("RemoveOldFilesWorker:", err)
			return
		}
	}
}

func (s *Server) RemoveJunkWorker(files []db.GetLastSeenAllRow, uploadDirectory string) {
	fsFiles, err := os.ReadDir(uploadDirectory)
	if err != nil {
		log.Println("RemoveJunkWorker:", err)
		return
	}
	found := false
	for _, direntry := range fsFiles {
		for _, file := range files {
			found = file.FileUuid == direntry.Name()
			if found {
				break
			}
		}
		if !found {
			err = os.RemoveAll(path.Join(uploadDirectory, direntry.Name()))
			if err != nil {
				log.Println("RemoveJunkWorker:", err)
				return
			}
		}
	}
}
