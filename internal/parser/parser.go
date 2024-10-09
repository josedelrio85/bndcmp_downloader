package parser

import (
	"io"

	"golang.org/x/net/html"
)

type parseClient struct{}

func NewParseClient() *parseClient {
	return &parseClient{}
}

func (p *parseClient) Parse(data io.Reader) (*html.Node, error) {
	node, err := html.Parse(data)
	if err != nil {
		return nil, err
	}
	return node, nil
}
