package scrapper

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/josedelrio85/bndcmp_downloader/internal/bandcamp"
	"github.com/josedelrio85/bndcmp_downloader/internal/model"
	"golang.org/x/net/html"
)

type TrackScrapper struct {
	URL        string
	Track      *model.Track
	httpClient Retriever
}

func NewTrackScrapper(URL string, httpClient Retriever) *TrackScrapper {
	return &TrackScrapper{
		URL:        URL,
		Track:      &model.Track{},
		httpClient: httpClient,
	}
}

func (t *TrackScrapper) Retrieve(url string) (io.Reader, error) {
	return t.httpClient.Retrieve(url)
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

func (t *TrackScrapper) Execute() error {
	fmt.Println("starting track scrapper for url: ", t.URL)
	reader, err := t.Retrieve(t.URL)
	if err != nil {
		fmt.Println("get ", err)
		panic(err)
	}

	node, err := t.Parse(reader)
	if err != nil {
		fmt.Println("parse ", err)
		panic(err)
	}

	err = t.Find(node)
	if err != nil {
		fmt.Println("find ", err)
		panic(err)
	}

	if t.Track != nil {
		fmt.Println("processing download for track: ", t.Track.Title)
		if t.Track.DownloadURL != "" {
			mp3_reader, err := t.Retrieve(t.Track.DownloadURL)
			if err != nil {
				fmt.Println("mp3 get ", err)
				panic(err)
			}

			trackTitle := t.Track.Title + ".mp3"
			if err := t.Save(mp3_reader, trackTitle, nil); err != nil {
				fmt.Println("save ", err)
				panic(err)
			}
		}
	}

	return nil
}
