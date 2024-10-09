package scrapper

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/josedelrio85/bndcmp_downloader/internal/bandcamp"
	"github.com/josedelrio85/bndcmp_downloader/internal/model"
	"golang.org/x/net/html"
)

type TrackScrapper struct {
	Track *model.Track
}

func NewTrackScrapper() *TrackScrapper {
	return &TrackScrapper{
		Track: &model.Track{},
	}
}

func (t *TrackScrapper) Retrieve(url string) (io.Reader, error) {
	return retrieve(url)
}

func retrieve(url string) (io.Reader, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func (t *TrackScrapper) Parse(data io.Reader) (*html.Node, error) {
	return parse(data)
}

func parse(data io.Reader) (*html.Node, error) {
	node, err := html.Parse(data)
	if err != nil {
		return nil, err
	}
	return node, nil
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

func (t *TrackScrapper) Save(data io.Reader, filename string, folder *string) error {
	return save(data, filename, folder)
}

func save(data io.Reader, filename string, folder *string) error {
	folderPath := "downloads"
	if folder != nil {
		folderPath = *folder
	}

	newFile, err := os.Create(folderPath + "/" + filename)
	if err != nil {
		return err
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, data)
	if err != nil {
		return err
	}

	return nil
}
