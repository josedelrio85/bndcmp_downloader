package scrapper

import (
	"io"

	"golang.org/x/net/html"
)

//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=mock_$GOFILE
type Scrapper interface {
	Retriever
	Parser
	Finder
	Saver
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
	Save(data io.Reader, filename string, folder *string) error
}
