package scrapper

import (
	io "io"
	"log"
	"net/url"
	"regexp"
	"strings"

	html "golang.org/x/net/html"
)

type AlbumScrapper struct {
	URL              *url.URL
	TrackList        []string
	httpClient       Retriever
	parseClient      Parser
	saveClient       Saver
	executeClient    func(string, Retriever, Parser, Saver, *map[string]bool) Executer
	downloadedTracks *map[string]bool
}

func NewAlbumScrapper(url *url.URL, httpClient Retriever, parseClient Parser, saveClient Saver, downloadedTracks *map[string]bool) *AlbumScrapper {
	return &AlbumScrapper{
		URL:         url,
		TrackList:   []string{},
		httpClient:  httpClient,
		parseClient: parseClient,
		saveClient:  saveClient,
		executeClient: func(url string, httpClient Retriever, parseClient Parser, saveClient Saver, downloadedTracks *map[string]bool) Executer {
			return NewTrackScrapper(url, httpClient, parseClient, saveClient, downloadedTracks)
		},
		downloadedTracks: downloadedTracks,
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
	log.Println("Scrapping album at:", a.URL.String())
	reader, err := a.Retrieve(a.URL.String())
	if err != nil {
		log.Println("Error retrieving album:", err)
		return err
	}

	node, err := a.Parse(reader)
	if err != nil {
		log.Println("Error parsing album HTML:", err)
		return err
	}

	err = a.Find(node)
	if err != nil {
		log.Println("Error finding tracks in album HTML:", err)
		return err
	}

	baseURL := url.URL{
		Scheme: a.URL.Scheme,
		Host:   a.URL.Host,
	}
	log.Printf("%d tracks to download \n", len(a.TrackList))
	for _, track := range a.TrackList {
		trackURL := baseURL.ResolveReference(&url.URL{Path: track})
		log.Println("Retrieving track:", trackURL.String())
		trackScrapper := a.executeClient(trackURL.String(), a.httpClient, a.parseClient, a.saveClient, a.downloadedTracks)
		if err := trackScrapper.Execute(); err != nil {
			log.Println("Error executing track scrapper:", err)
			return err
		}
	}
	return nil
}
