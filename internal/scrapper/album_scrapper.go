package scrapper

import (
	"fmt"
	io "io"
	"net/url"
	"regexp"
	"strings"

	html "golang.org/x/net/html"
)

type AlbumScrapper struct {
	URL           *url.URL
	TrackList     []string
	httpClient    Retriever
	parseClient   Parser
	saveClient    Saver
	executeClient func(string, Retriever, Parser, Saver) Executer
}

func NewAlbumScrapper(url *url.URL, httpClient Retriever, parseClient Parser, saveClient Saver) *AlbumScrapper {
	return &AlbumScrapper{
		URL:         url,
		TrackList:   []string{},
		httpClient:  httpClient,
		parseClient: parseClient,
		saveClient:  saveClient,
		executeClient: func(url string, httpClient Retriever, parseClient Parser, saveClient Saver) Executer {
			return NewTrackScrapper(url, httpClient, parseClient, saveClient)
		},
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

func (a *AlbumScrapper) Save(data io.Reader, filename string) error {
	return nil
}

func (a *AlbumScrapper) Execute() error {
	fmt.Println("starting album scrapper for url: ", a.URL.String())
	reader, err := a.Retrieve(a.URL.String())
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
		Scheme: a.URL.Scheme,
		Host:   a.URL.Host,
	}
	for _, track := range a.TrackList {
		trackURL := baseURL.ResolveReference(&url.URL{Path: track})
		fmt.Println("retrieving track: ", trackURL.String())
		trackScrapper := a.executeClient(trackURL.String(), a.httpClient, a.parseClient, a.saveClient)
		if err := trackScrapper.Execute(); err != nil {
			fmt.Println("execute ", err)
			return err
		}
	}
	return nil
}
