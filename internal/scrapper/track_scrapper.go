package scrapper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/josedelrio85/bndcmp_downloader/internal/album_catalog"
	"github.com/josedelrio85/bndcmp_downloader/internal/bandcamp"
	"github.com/josedelrio85/bndcmp_downloader/internal/model"
	"golang.org/x/net/html"
)

type TrackScrapper struct {
	Track        *model.Track
	httpClient   Retriever
	parseClient  Parser
	saveClient   Saver
	albumCatalog album_catalog.AlbumCatalog
}

func NewTrackScrapper(httpClient Retriever, parseClient Parser, saveClient Saver, albumCatalog album_catalog.AlbumCatalog) *TrackScrapper {
	return &TrackScrapper{
		Track:        &model.Track{},
		httpClient:   httpClient,
		parseClient:  parseClient,
		saveClient:   saveClient,
		albumCatalog: albumCatalog,
	}
}

func (t *TrackScrapper) Retrieve(url string) (io.Reader, error) {
	return t.httpClient.Retrieve(url)
}

func (t *TrackScrapper) Parse(data io.Reader) (*html.Node, error) {
	return t.parseClient.Parse(data)
}

func (t *TrackScrapper) Find(node *html.Node) error {
	if node.Type == html.ElementNode && node.Data == "script" {
		for _, z := range node.Attr {
			if z.Key == "data-tralbum" {
				var albumInfo bandcamp.TrAlbum
				if err := json.Unmarshal([]byte(z.Val), &albumInfo); err != nil {
					return err
				}
				t.Track = albumInfo.ToTrack()
				return nil
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if err := t.Find(c); err != nil {
			return err
		}
	}
	return nil
}

func (t *TrackScrapper) Save(data io.Reader, track *model.Track) error {
	return t.saveClient.Save(data, track)
}

func (t *TrackScrapper) Execute(trackURL *url.URL) error {
	log.Printf("Starting track scrapper for URL: %s", trackURL.String())
	reader, err := t.Retrieve(trackURL.String())
	if err != nil {
		log.Printf("Error retrieving URL %s: %v", trackURL.String(), err)
		return err
	}

	node, err := t.Parse(reader)
	if err != nil {
		log.Printf("Error parsing HTML content: %v", err)
		return err
	}

	err = t.Find(node)
	if err != nil {
		log.Printf("Error finding track information in HTML: %v", err)
		return err
	}

	if t.Track != nil {
		if t.isDownloaded() {
			return nil
		}
		log.Printf("Processing download for track: %s", t.Track.Title)
		if t.Track.DownloadURL != "" {
			mp3_reader, err := t.Retrieve(t.Track.DownloadURL)
			if err != nil {
				log.Printf("Error retrieving MP3 from URL %s: %v", t.Track.DownloadURL, err)
				return err
			}

			if err := t.Save(mp3_reader, t.Track); err != nil {
				log.Printf("Error saving track %s: %v", t.Track.Title, err)
				return err
			}

			t.updateDownloadedTracks()
		}
	}

	return nil
}

func (t *TrackScrapper) isDownloaded() bool {
	mapDir := t.albumCatalog.GetMapDir()
	if mapDir != nil {
		filePath := t.generateFilePath()
		log.Printf("Checking if track %s is downloaded", filePath)
		if _, ok := (*mapDir)[filePath]; ok {
			log.Printf("Track %s already downloaded", filePath)
			return true
		}
	}
	return false
}

func (t *TrackScrapper) generateFilePath() string {
	trackName := fmt.Sprintf("%02d - %s.mp3", t.Track.TrackNumber, t.Track.Title)
	if t.Track.Album == nil {
		return strings.Join([]string{t.Track.Artist, trackName}, "/")
	}
	return strings.Join([]string{t.Track.Artist, *t.Track.Album, trackName}, "/")
}

func (t *TrackScrapper) updateDownloadedTracks() {
	if t.albumCatalog.GetMapDir() != nil {
		filePath := t.generateFilePath()
		t.albumCatalog.Update(filePath)
	}
}
