package utils

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

// Check ifthe given fileName is a zip file
// Only need to read first 512 bytes @see http://golang.org/pkg/net/http/#DetectContentType
func IsZip(fileName string) bool {

	file, _ := os.Open(fileName)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Errorf("unable to close file %v due to %v", fileName, err)
		}
	}(file)

	buff := make([]byte, 512)

	_, err := file.Read(buff)
	if err != nil {
		log.Errorf("unable to check if file is a zip due to %v", err)
		return false
	}

	fileType := http.DetectContentType(buff)

	switch fileType {
	case "application/x-gzip", "application/zip":
		return true
	default:
		return false
	}
}
