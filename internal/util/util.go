package util

import (
	"log"
	"os"
	"path/filepath"
)

var MapDir = make(map[string]bool)
var initialBaseFolder string

func RecursivelyListDirectory(baseFolder string) {
	if initialBaseFolder == "" {
		initialBaseFolder = baseFolder
	}

	entries, err := os.ReadDir(baseFolder)
	if err != nil {
		log.Println("error iterating over folder: ", err)
		MapDir = make(map[string]bool)
		return
	}

	for _, entry := range entries {
		completePath := filepath.Join(baseFolder, entry.Name())
		relativePath, err := filepath.Rel(initialBaseFolder, completePath)
		if err != nil {
			log.Println("error getting relative path:", err)
			MapDir = make(map[string]bool)
			return
		}
		MapDir[relativePath] = true
		if entry.IsDir() {
			RecursivelyListDirectory(completePath)
		}
	}
}
