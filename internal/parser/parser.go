package parser

import (
	"io"

	"golang.org/x/net/html"
)

type ParseClient struct{}

func NewParseClient() *ParseClient {
	return &ParseClient{}
}

func (p *ParseClient) Parse(data io.Reader) (*html.Node, error) {
	node, err := html.Parse(data)
	if err != nil {
		return nil, err
	}
	return node, nil
}
