package bandcamp

import (
	"errors"
	"net/url"
	"strings"
)

type BandcampURL struct {
	Value string
	URL   *url.URL
}

func (b *BandcampURL) Parse() error {
	parsedURL, err := url.Parse(b.Value)
	if err != nil {
		return err
	}

	b.URL = parsedURL

	return nil
}

func (b *BandcampURL) Validate() error {
	if !strings.Contains(b.URL.Host, "bandcamp.com") {
		return errors.New("invalid Bandcamp url")
	}

	return nil
}

// Classify returns the type of Bandcamp URL
func (b *BandcampURL) Classify() URLType {
	// Remove leading slash and split the path
	path := strings.TrimPrefix(b.URL.Path, "/")
	if path == "" {
		return URLTypeUnknown
	}

	// Get the first segment of the path
	segment := strings.SplitN(path, "/", 2)[0]

	switch segment {
	case "track":
		return URLTypeTrack
	case "album":
		return URLTypeAlbum
	case "music":
		return URLTypeDiscography
	default:
		return URLTypeUnknown
	}
}

// URLType represents the type of Bandcamp URL
type URLType int

const (
	URLTypeUnknown URLType = iota
	URLTypeTrack
	URLTypeAlbum
	URLTypeDiscography
)
