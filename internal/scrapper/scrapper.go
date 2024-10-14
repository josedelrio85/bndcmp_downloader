package scrapper

import (
	"io"
	"net/url"

	"github.com/josedelrio85/bndcmp_downloader/internal/model"
	"golang.org/x/net/html"
)

//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=mock_$GOFILE
type Scrapper interface {
	Retriever
	Parser
	Finder
	Saver
	Executer
}

type Retriever interface {
	Retrieve(url string) (io.Reader, error)
}

type Parser interface {
	Parse(data io.Reader) (*html.Node, error)
}

type Finder interface {
	Find(node *html.Node) error
}

type Saver interface {
	Save(data io.Reader, track *model.Track) error
}

type Executer interface {
	Execute(resourceURL *url.URL) error
}
