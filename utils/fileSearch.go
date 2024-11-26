package utils

import (
	"os"
	"path/filepath"
)

func FileHandle(filePath string) (*os.File, error) {
	filePath, err := filepath.Rel("/OnlineMusic", filePath)
	if err := IsDirExists(filePath); err != nil {
		return &os.File{}, err
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return &os.File{}, err
	}
	return file, nil
}

func IsDirExists(filePath string) error {
	dirPath := filepath.Dir(filePath)
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return err
	}
	return nil
}
