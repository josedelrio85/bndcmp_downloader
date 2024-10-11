package main

import (
	"log"

	"github.com/josedelrio85/bndcmp_downloader/internal/parser"
	"github.com/josedelrio85/bndcmp_downloader/internal/prompt"
	"github.com/josedelrio85/bndcmp_downloader/internal/retriever"
	"github.com/josedelrio85/bndcmp_downloader/internal/saver"
	"github.com/josedelrio85/bndcmp_downloader/internal/scrapper"
	"github.com/josedelrio85/bndcmp_downloader/internal/util"
)

func main() {
	log.Println("Starting Bandcamp downloader")

	promptChain := setupPromptChain()

	httpClient, parseClient, saveClient := setup(&promptChain.ChainMessage.StorageType)

	util.RecursivelyListDirectory(promptChain.ChainMessage.StorageType)

	var err error
	switch promptChain.ChainMessage.ScrapType {
	case scrapper.Track:
		err = scrapper.NewTrackScrapper(promptChain.ChainMessage.URL.Value, httpClient, parseClient, saveClient, &util.MapDir).Execute()
	case scrapper.Album:
		err = scrapper.NewAlbumScrapper(promptChain.ChainMessage.URL.URL, httpClient, parseClient, saveClient, &util.MapDir).Execute()
	case scrapper.Discography:
		err = scrapper.NewDiscographyScrapper(promptChain.ChainMessage.URL.URL, httpClient, parseClient, saveClient, &util.MapDir).Execute()
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
