package setup

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/josedelrio85/bndcmp_downloader/internal/album_catalog"
	"github.com/josedelrio85/bndcmp_downloader/internal/parser"
	"github.com/josedelrio85/bndcmp_downloader/internal/retriever"
	"github.com/josedelrio85/bndcmp_downloader/internal/saver"
)

type Config struct {
	BaseFolder   string
	Retriever    *retriever.HttpClient
	Parser       *parser.ParseClient
	Saver        *saver.LocalSaver
	AlbumCatalog album_catalog.AlbumCatalog
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	baseFolder := loadBaseFolder()

	albumCatalog := album_catalog.NewInMemoryAlbumCatalog(baseFolder)
	if err := albumCatalog.Generate(baseFolder); err != nil {
		log.Fatal("Error generating album catalog: ", err)
	}

	return &Config{
		BaseFolder:   baseFolder,
		Retriever:    retriever.NewHttpClient(),
		Parser:       parser.NewParseClient(),
		Saver:        saver.NewLocalSaver(&baseFolder),
		AlbumCatalog: albumCatalog,
	}
}

func loadBaseFolder() string {
	baseFolder := os.Getenv("BASE_FOLDER")
	if baseFolder == "" {
		log.Fatal("BASE_FOLDER is not set")
	}

	return baseFolder
}
