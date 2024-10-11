package main

import (
	"fmt"
	"net/url"

	"github.com/josedelrio85/bndcmp_downloader/internal/parser"
	"github.com/josedelrio85/bndcmp_downloader/internal/retriever"
	"github.com/josedelrio85/bndcmp_downloader/internal/saver"
	"github.com/josedelrio85/bndcmp_downloader/internal/scrapper"
)

func main() {
	fmt.Println("starting Bandcamp downloader")

	httpClient := retriever.NewHttpClient()
	parseClient := parser.NewParseClient()
	saveFolder := "downloads"
	saveClient := saver.NewLocalSaver(&saveFolder)

	// trackURL := "https://kinggizzard.bandcamp.com/track/elbow"
	// trackScrapper := scrapper.NewTrackScrapper(trackURL, httpClient, parseClient, saveClient)
	// if err := trackScrapper.Execute(); err != nil {
	// 	fmt.Println("execute ", err)
	// 	return
	// }

	// albumURL := "https://kinggizzard.bandcamp.com/album/12-bar-bruise"
	// parsedAlbumURL, err := url.Parse(albumURL)
	// if err != nil {
	// 	fmt.Println("parsing album url ", err)
	// 	return
	// }
	// albumScrapper := scrapper.NewAlbumScrapper(parsedAlbumURL, httpClient, parseClient, saveClient)
	// if err := albumScrapper.Execute(); err != nil {
	// 	fmt.Println("execute ", err)
	// 	return
	// }

	discographyURL := "https://kinggizzard.bandcamp.com/music"
	parsedDiscographyURL, err := url.Parse(discographyURL)
	if err != nil {
		fmt.Println("parsing discography url ", err)
		return
	}
	discographyScrapper := scrapper.NewDiscographyScrapper(parsedDiscographyURL, httpClient, parseClient, saveClient)
	if err := discographyScrapper.Execute(); err != nil {
		fmt.Println("execute ", err)
		return
	}

}
