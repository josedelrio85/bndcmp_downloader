package saver

import (
	"io"
	"os"
	"path/filepath"
)

type LocalSaver struct {
	folder *string
}

func NewLocalSaver(folder *string) *LocalSaver {
	return &LocalSaver{folder: folder}
}

func (s *LocalSaver) Save(data io.Reader, filename string) error {
	folder := ""
	if s.folder != nil {
		folder = *s.folder
	}
	if folder != "" {
		// Check if the folder exists
		info, err := os.Stat(folder)
		if os.IsNotExist(err) {
			return err
		}

		if !info.IsDir() {
			return os.ErrNotExist
		}
	}

	path := filepath.Join(folder, filename)
	newFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, data)
	if err != nil {
		return err
	}

	return nil
}
