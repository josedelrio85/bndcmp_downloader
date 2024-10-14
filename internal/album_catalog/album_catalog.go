package album_catalog

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=mock_$GOFILE
type AlbumCatalog interface {
	Generate(folder string) error
	GetMapDir() *map[string]bool
	Update(path string)
}

type InMemoryAlbumCatalog struct {
	mapDir     map[string]bool
	baseFolder string
	mutex      sync.Mutex
}

func NewInMemoryAlbumCatalog(baseFolder string) *InMemoryAlbumCatalog {
	return &InMemoryAlbumCatalog{
		mapDir:     make(map[string]bool),
		baseFolder: baseFolder,
		mutex:      sync.Mutex{},
	}
}

func (i *InMemoryAlbumCatalog) Generate(folder string) error {
	if i.baseFolder == "" {
		i.baseFolder = folder
	}
	entries, err := os.ReadDir(folder)
	if err != nil {
		log.Printf("error iterating over folder: %s: %v\n", folder, err)
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			next := filepath.Join(folder, entry.Name())
			i.Generate(next)
		} else {
			i.mutex.Lock()
			nextTrack := filepath.Join(folder, entry.Name())
			nextTrack = strings.TrimPrefix(nextTrack, i.baseFolder)
			nextTrack = strings.TrimPrefix(nextTrack, string(os.PathSeparator))
			i.mapDir[nextTrack] = true
			i.mutex.Unlock()
		}
	}
	return nil
}

func (i *InMemoryAlbumCatalog) GetMapDir() *map[string]bool {
	return &i.mapDir
}

func (i *InMemoryAlbumCatalog) Update(path string) {
	i.mutex.Lock()
	i.mapDir[path] = true
	i.mutex.Unlock()
}
