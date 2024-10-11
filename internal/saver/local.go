package saver

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/josedelrio85/bndcmp_downloader/internal/model"
)

type LocalSaver struct {
	storageFolder string
}

func NewLocalSaver(folder *string) *LocalSaver {
	storageFolder := "./"
	if folder != nil {
		storageFolder = *folder
	}
	storageFolder = strings.TrimSuffix(storageFolder, "/")
	return &LocalSaver{storageFolder: storageFolder}
}

func (s *LocalSaver) Save(data io.Reader, track *model.Track) error {
	if track == nil {
		return errors.New("track is nil")
	}

	directoryStructure := s.generateDirectoryStructure(track)
	directoryStructureWithBase := filepath.Join(s.storageFolder, directoryStructure)
	if err := s.checkFolder(directoryStructureWithBase); err != nil {
		return err
	}

	trackName := fmt.Sprintf("%02d - %s.mp3", track.TrackNumber, track.Title)
	if err := s.saveFile(directoryStructureWithBase, trackName, data); err != nil {
		return err
	}
	return nil
}

func (s *LocalSaver) generateDirectoryStructure(track *model.Track) string {
	if track.Album == nil {
		return strings.Join([]string{track.Artist}, "/")
	}
	return strings.Join([]string{track.Artist, *track.Album}, "/")
}

func (s *LocalSaver) checkFolder(base string) error {
	_, err := os.Stat(base)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(base, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func (s *LocalSaver) saveFile(base string, filename string, data io.Reader) error {
	filePath := filepath.Join(base, filename)
	fmt.Println(filePath)
	newFile, err := os.Create(filePath)
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
