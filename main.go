package main

import (
	"fmt"

	"github.com/josedelrio85/bndcmp_downloader/internal/scrapper"
)

func main() {
	fmt.Println("starting Bandcamp downloader")

	// elbowURL := "https://kinggizzard.bandcamp.com/track/elbow"
	// trackScrapper(elbowURL)

	albumURL := "https://kinggizzard.bandcamp.com/album/12-bar-bruise"
	albumScrapper(albumURL)
}

func trackScrapper(url string) {
	fmt.Println("starting track scrapper for url: ", url)
	trackScrapper := scrapper.NewTrackScrapper()
	reader, err := trackScrapper.Retrieve(url)
	if err != nil {
		fmt.Println("get ", err)
		panic(err)
	}

	node, err := trackScrapper.Parse(reader)
	if err != nil {
		fmt.Println("parse ", err)
		panic(err)
	}

	err = trackScrapper.Find(node)
	if err != nil {
		fmt.Println("find ", err)
		panic(err)
	}

	if trackScrapper.Track != nil {
		fmt.Println("processing download for track: ", trackScrapper.Track.Title)
		if trackScrapper.Track.DownloadURL != "" {
			mp3_reader, err := trackScrapper.Retrieve(trackScrapper.Track.DownloadURL)
			if err != nil {
				fmt.Println("mp3 get ", err)
				panic(err)
			}

			trackTitle := trackScrapper.Track.Title + ".mp3"
			if err := trackScrapper.Save(mp3_reader, trackTitle, nil); err != nil {
				fmt.Println("save ", err)
				panic(err)
			}
		}
	}
}

func albumScrapper(url string) {
	fmt.Println("starting album scrapper for url: ", url)
	albumScrapper := scrapper.NewAlbumScrapper()
	reader, err := albumScrapper.Retrieve(url)
	if err != nil {
		fmt.Println("get ", err)
		panic(err)
	}

	node, err := albumScrapper.Parse(reader)
	if err != nil {
		fmt.Println("parse ", err)
		panic(err)
	}

	err = albumScrapper.Find(node)
	if err != nil {
		fmt.Println("find ", err)
		panic(err)
	}

	for _, track := range albumScrapper.TrackList {
		trackURL := "https://kinggizzard.bandcamp.com" + track
		fmt.Println("retrieving track: ", trackURL)
		trackScrapper(trackURL)
	}
}
