package scrapper

import (
	io "io"
	"regexp"
	"strings"

	html "golang.org/x/net/html"
)

type AlbumScrapper struct {
	TrackList   []string
	httpClient  Retriever
	parseClient Parser
}

func NewAlbumScrapper(httpClient Retriever, parseClient Parser) *AlbumScrapper {
	return &AlbumScrapper{
		TrackList:   []string{},
		httpClient:  httpClient,
		parseClient: parseClient,
	}
}

func (a *AlbumScrapper) Retrieve(url string) (io.Reader, error) {
	return a.httpClient.Retrieve(url)
}

func (a *AlbumScrapper) Parse(data io.Reader) (*html.Node, error) {
	return a.parseClient.Parse(data)
}

func (a *AlbumScrapper) Find(node *html.Node) error {
	if err := a.find(node); err != nil {
		return err
	}
	a.TrackList = a.processTrackList()
	return nil
}

func (a *AlbumScrapper) find(node *html.Node) error {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" && strings.Contains(attr.Val, "track") {
				a.TrackList = append(a.TrackList, attr.Val)
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

func (a *AlbumScrapper) processTrackList() []string {
	pattern := regexp.MustCompile(`^/track/[a-z0-9-]+$`)
	seen := make(map[string]bool)
	var result []string

	for _, track := range a.TrackList {
		if pattern.MatchString(track) && !seen[track] {
			result = append(result, track)
			seen[track] = true
		}
	}

	return result
}

func (a *AlbumScrapper) Save(data io.Reader, filename string, folder *string) error {
	return save(data, filename, folder)
}
