package scrapper

import (
	"encoding/json"
	"fmt"
	io "io"
	"log"
	"net/url"
	"regexp"

	"github.com/josedelrio85/bndcmp_downloader/internal/bandcamp"
	html "golang.org/x/net/html"
)

type DiscographyScrapper struct {
	discographyURL   *url.URL
	AlbumList        []string
	httpClient       Retriever
	parseClient      Parser
	saveClient       Saver
	executeClient    func(*url.URL, Retriever, Parser, Saver, *map[string]bool) Executer
	downloadedTracks *map[string]bool
}

func NewDiscographyScrapper(discographyURL *url.URL, httpClient Retriever, parseClient Parser, saveClient Saver, downloadedTracks *map[string]bool) *DiscographyScrapper {
	return &DiscographyScrapper{
		discographyURL: discographyURL,
		AlbumList:      []string{},
		httpClient:     httpClient,
		parseClient:    parseClient,
		saveClient:     saveClient,
		executeClient: func(url *url.URL, httpClient Retriever, parseClient Parser, saveClient Saver, downloadedTracks *map[string]bool) Executer {
			return NewAlbumScrapper(url, httpClient, parseClient, saveClient, downloadedTracks)
		},
		downloadedTracks: downloadedTracks,
	}
}

func (a *DiscographyScrapper) Retrieve(url string) (io.Reader, error) {
	return a.httpClient.Retrieve(url)
}

func (a *DiscographyScrapper) Parse(data io.Reader) (*html.Node, error) {
	return a.parseClient.Parse(data)
}

func (a *DiscographyScrapper) Find(node *html.Node) error {
	if err := a.find(node); err != nil {
		return err
	}
	a.AlbumList = a.processAlbumList()
	return nil
}

func (a *DiscographyScrapper) find(node *html.Node) error {
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
		if err := a.find(c); err != nil {
			return err
		}
	}
	return nil
}

func (a *DiscographyScrapper) processAlbumList() []string {
	pattern := regexp.MustCompile(`^/album/[a-z0-9-]+$`)
	seen := make(map[string]bool)
	var result []string

	fmt.Println("processAlbumList len albumlist ", len(a.AlbumList))
	for _, track := range a.AlbumList {
		if pattern.MatchString(track) && !seen[track] {
			result = append(result, track)
			seen[track] = true
		}
	}

	return result
}

func (a *DiscographyScrapper) Save(data io.Reader, filename string) error {
	return nil
}

func (a *DiscographyScrapper) Execute() error {
	log.Printf("Starting discography scrapper for URL: %s", a.discographyURL.String())
	reader, err := a.Retrieve(a.discographyURL.String())
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
		Scheme: a.discographyURL.Scheme,
		Host:   a.discographyURL.Host,
	}
	log.Printf("%d albums to download \n", len(a.AlbumList))
	for _, album := range a.AlbumList {
		albumURL := baseURL.ResolveReference(&url.URL{Path: album})
		log.Printf("Retrieving album: %s", albumURL.String())
		albumScrapper := a.executeClient(albumURL, a.httpClient, a.parseClient, a.saveClient, a.downloadedTracks)
		if err := albumScrapper.Execute(); err != nil {
			log.Printf("Error executing album scrapper for %s: %v", albumURL.String(), err)
			return err
		}
	}
	return nil
}
