package utils

import "os"

func FileSize(relativePath string) (int64, error) {
	filestat, err := os.Stat(relativePath)
	if err != nil {
		return 0, err
	}
	return filestat.Size(), nil
}

func FileSizeGuarded(relativePath string) int64 {
	fs, err := FileSize(relativePath)
	Check(err)
	return fs
}
