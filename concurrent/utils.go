package concurrent

import (
	"os"
	"path/filepath"
)

func WriteFile(input string, filePath string) error {
	data := []byte(input)
	dir, _ := filepath.Split(filePath)

	if _, err := os.Stat(dir); err == nil {
		os.Remove(filePath)
	} else {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filePath, data, 0); err != nil {
		return err
	}
	return nil
}
