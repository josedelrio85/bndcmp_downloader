package scrapper

import (
	"fmt"
	io "io"
	"net/url"
	"regexp"
	"strings"

	html "golang.org/x/net/html"
)

type DiscographyScrapper struct {
	discographyURL *url.URL
	AlbumList      []string
	httpClient     Retriever
	parseClient    Parser
	saveClient     Saver
	executeClient  func(*url.URL, Retriever, Parser, Saver) Executer
}

func NewDiscographyScrapper(discographyURL *url.URL, httpClient Retriever, parseClient Parser, saveClient Saver) *DiscographyScrapper {
	return &DiscographyScrapper{
		discographyURL: discographyURL,
		AlbumList:      []string{},
		httpClient:     httpClient,
		parseClient:    parseClient,
		saveClient:     saveClient,
		executeClient: func(url *url.URL, httpClient Retriever, parseClient Parser, saveClient Saver) Executer {
			return NewAlbumScrapper(url, httpClient, parseClient, saveClient)
		},
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
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" && strings.Contains(attr.Val, "album") {
				a.AlbumList = append(a.AlbumList, attr.Val)
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
	fmt.Println("starting discography scrapper for url: ", a.discographyURL.String())
	reader, err := a.Retrieve(a.discographyURL.String())
	if err != nil {
		fmt.Println("get ", err)
		return err
	}

	node, err := a.Parse(reader)
	if err != nil {
		fmt.Println("parse ", err)
		return err
	}

	err = a.Find(node)
	if err != nil {
		fmt.Println("find ", err)
		return err
	}

	baseURL := url.URL{
		Scheme: a.discographyURL.Scheme,
		Host:   a.discographyURL.Host,
	}
	for _, album := range a.AlbumList {
		albumURL := baseURL.ResolveReference(&url.URL{Path: album})
		fmt.Println("retrieving album: ", albumURL.String())
		albumScrapper := a.executeClient(albumURL, a.httpClient, a.parseClient, a.saveClient)
		if err := albumScrapper.Execute(); err != nil {
			fmt.Println("execute ", err)
			return err
		}
	}
	return nil
}
