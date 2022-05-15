package util

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

func FileExists(path string) bool {
	_, statErr := os.Stat(path)
	return !os.IsNotExist(statErr)
}

func GetFile(filePath string) (file *os.File, err error) {
	dir := filepath.Dir(filePath)

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		errDir := os.MkdirAll(dir, 0700)
		if errDir != nil {
			return nil, errDir
		}
	}

	return os.Create(filePath)
}

func GetHash(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
