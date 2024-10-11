package scrapper

import (
	"encoding/json"
	"io"
	"log"

	"github.com/josedelrio85/bndcmp_downloader/internal/bandcamp"
	"github.com/josedelrio85/bndcmp_downloader/internal/model"
	"golang.org/x/net/html"
)

type TrackScrapper struct {
	URL         string
	Track       *model.Track
	httpClient  Retriever
	parseClient Parser
	saveClient  Saver
}

func NewTrackScrapper(URL string, httpClient Retriever, parseClient Parser, saveClient Saver) *TrackScrapper {
	return &TrackScrapper{
		URL:         URL,
		Track:       &model.Track{},
		httpClient:  httpClient,
		parseClient: parseClient,
		saveClient:  saveClient,
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

func (t *TrackScrapper) Execute() error {
	log.Printf("Starting track scrapper for URL: %s", t.URL)
	reader, err := t.Retrieve(t.URL)
	if err != nil {
		log.Printf("Error retrieving URL %s: %v", t.URL, err)
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
		}
	}

	return nil
}
