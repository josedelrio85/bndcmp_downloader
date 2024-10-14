package album_catalog

import (
	"log"
	"os"
	"path/filepath"
)

type AlbumCatalog interface {
	Generate() error
}

type InMemoryAlbumCatalog struct {
	MapDir     map[string]bool
	baseFolder string
}

func NewInMemoryAlbumCatalog(baseFolder string) *InMemoryAlbumCatalog {
	return &InMemoryAlbumCatalog{
		MapDir:     make(map[string]bool),
		baseFolder: baseFolder,
	}
}

func (i *InMemoryAlbumCatalog) Generate() error {
	entries, err := os.ReadDir(i.baseFolder)
	if err != nil {
		log.Println("error iterating over folder: ", err)
		return err
	}

	for _, entry := range entries {
		completePath := filepath.Join(i.baseFolder, entry.Name())
		relativePath, err := filepath.Rel(i.baseFolder, completePath)
		if err != nil {
			log.Println("error getting relative path:", err)
			return err
		}
		i.MapDir[relativePath] = true
		if entry.IsDir() {
			i.baseFolder = completePath
			i.Generate()
		}
	}

	return nil
}
