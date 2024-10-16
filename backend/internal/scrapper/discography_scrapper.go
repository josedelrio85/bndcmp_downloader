package scrapper

import (
	"encoding/json"
	io "io"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/josedelrio85/bndcmp_downloader/internal/album_catalog"
	"github.com/josedelrio85/bndcmp_downloader/internal/bandcamp"
	model "github.com/josedelrio85/bndcmp_downloader/internal/model"
	html "golang.org/x/net/html"
)

type DiscographyScrapper struct {
	AlbumList     []string
	httpClient    Retriever
	parseClient   Parser
	saveClient    Saver
	executeClient func(Retriever, Parser, Saver, album_catalog.AlbumCatalog) Executer
	albumCatalog  album_catalog.AlbumCatalog
}

func NewDiscographyScrapper(httpClient Retriever, parseClient Parser, saveClient Saver, albumCatalog album_catalog.AlbumCatalog) *DiscographyScrapper {
	return &DiscographyScrapper{
		AlbumList:   []string{},
		httpClient:  httpClient,
		parseClient: parseClient,
		saveClient:  saveClient,
		executeClient: func(httpClient Retriever, parseClient Parser, saveClient Saver, albumCatalog album_catalog.AlbumCatalog) Executer {
			return NewAlbumScrapper(httpClient, parseClient, saveClient, albumCatalog)
		},
		albumCatalog: albumCatalog,
	}
}

func (a *DiscographyScrapper) Retrieve(url string) (io.Reader, error) {
	return a.httpClient.Retrieve(url)
}

func (a *DiscographyScrapper) Parse(data io.Reader) (*html.Node, error) {
	return a.parseClient.Parse(data)
}

func (a *DiscographyScrapper) Find(node *html.Node) error {
	if err := a.findByDataClientItems(node); err != nil {
		return err
	}

	if len(a.AlbumList) == 0 {
		if err := a.findByHref(node); err != nil {
			return err
		}
	}

	a.AlbumList = a.processAlbumList()
	return nil
}

func (a *DiscographyScrapper) findByDataClientItems(node *html.Node) error {
	if node.Type == html.ElementNode && node.Data == "ol" {
		for _, attr := range node.Attr {
			if attr.Key == "data-client-items" {
				var bandcampAlbums []bandcamp.Album
				err := json.Unmarshal([]byte(attr.Val), &bandcampAlbums)
				if err != nil {
					log.Printf("Error unmarshalling Bandcamp albums: %v", err)
					return err
				}

				for _, album := range bandcampAlbums {
					a.AlbumList = append(a.AlbumList, album.PageURL)
				}
				return nil
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if err := a.findByDataClientItems(c); err != nil {
			return err
		}
	}
	return nil
}

func (a *DiscographyScrapper) findByHref(node *html.Node) error {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" && strings.Contains(attr.Val, "album") {
				a.AlbumList = append(a.AlbumList, attr.Val)
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if err := a.findByHref(c); err != nil {
			return err
		}
	}
	return nil
}

func (a *DiscographyScrapper) processAlbumList() []string {
	pattern := regexp.MustCompile(`^/album/[a-z0-9-]+$`)
	seen := make(map[string]bool)
	var result []string

	for _, track := range a.AlbumList {
		if pattern.MatchString(track) && !seen[track] {
			result = append(result, track)
			seen[track] = true
		}
	}

	return result
}

func (a *DiscographyScrapper) Save(data io.Reader, track *model.Track) error {
	return nil
}

func (a *DiscographyScrapper) Execute(discographyURL *url.URL) error {
	log.Printf("Starting discography scrapper for URL: %s", discographyURL.String())
	reader, err := a.Retrieve(discographyURL.String())
	if err != nil {
		log.Printf("Error retrieving discography page: %v", err)
		return err
	}

	node, err := a.Parse(reader)
	if err != nil {
		log.Printf("Error parsing discography HTML: %v", err)
		return err
	}

	err = a.Find(node)
	if err != nil {
		log.Printf("Error finding albums in discography HTML: %v", err)
		return err
	}

	baseURL := url.URL{
		Scheme: discographyURL.Scheme,
		Host:   discographyURL.Host,
	}
	log.Printf("%d albums to download \n", len(a.AlbumList))
	for _, album := range a.AlbumList {
		albumURL := baseURL.ResolveReference(&url.URL{Path: album})
		log.Printf("Retrieving album: %s", albumURL.String())
		albumScrapper := a.executeClient(a.httpClient, a.parseClient, a.saveClient, a.albumCatalog)
		if err := albumScrapper.Execute(albumURL); err != nil {
			log.Printf("Error executing album scrapper for %s: %v", albumURL.String(), err)
			return err
		}
	}
	return nil
}
