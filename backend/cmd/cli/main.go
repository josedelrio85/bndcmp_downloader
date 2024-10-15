package main

import (
	"log"

	"github.com/josedelrio85/bndcmp_downloader/internal/album_catalog"
	"github.com/josedelrio85/bndcmp_downloader/internal/parser"
	"github.com/josedelrio85/bndcmp_downloader/internal/prompt"
	"github.com/josedelrio85/bndcmp_downloader/internal/retriever"
	"github.com/josedelrio85/bndcmp_downloader/internal/saver"
	"github.com/josedelrio85/bndcmp_downloader/internal/scrapper"
)

func main() {
	log.Println("Starting Bandcamp downloader CLI")

	promptChain := setupPromptChain()

	httpClient, parseClient, saveClient := setup(&promptChain.ChainMessage.StorageType)

	inMemoryAlbumCatalog := album_catalog.NewInMemoryAlbumCatalog(promptChain.ChainMessage.StorageType)
	if err := inMemoryAlbumCatalog.Generate(promptChain.ChainMessage.StorageType); err != nil {
		log.Println("Error generating album catalog: ", err)
	}

	var err error
	switch promptChain.ChainMessage.ScrapType {
	case scrapper.Track:
		err = scrapper.NewTrackScrapper(httpClient, parseClient, saveClient, inMemoryAlbumCatalog).Execute(promptChain.ChainMessage.URL.URL)
	case scrapper.Album:
		err = scrapper.NewAlbumScrapper(httpClient, parseClient, saveClient, inMemoryAlbumCatalog).Execute(promptChain.ChainMessage.URL.URL)
	case scrapper.Discography:
		err = scrapper.NewDiscographyScrapper(httpClient, parseClient, saveClient, inMemoryAlbumCatalog).Execute(promptChain.ChainMessage.URL.URL)
	default:
		log.Println("Invalid scrap type")
	}

	if err != nil {
		log.Println("Error executing scrapper: ", err)
	}
}

func setup(saveFolder *string) (*retriever.HttpClient, *parser.ParseClient, *saver.LocalSaver) {
	httpClient := retriever.NewHttpClient()
	parseClient := parser.NewParseClient()
	saveClient := saver.NewLocalSaver(saveFolder)

	return httpClient, parseClient, saveClient
}

func setupPromptChain() *prompt.Chain {
	scrapTypeQuestionLink := prompt.NewScrapTypeQuestionLink()
	urlCheckerLink := prompt.NewURLCheckerLink()
	storageQuestionLink := prompt.NewStorageQuestionLink()
	links := []prompt.Link{scrapTypeQuestionLink, urlCheckerLink, storageQuestionLink}
	return prompt.NewChain(links)
}
