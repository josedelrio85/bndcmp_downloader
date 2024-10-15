package prompt

import (
	"github.com/josedelrio85/bndcmp_downloader/internal/bandcamp"
	"github.com/josedelrio85/bndcmp_downloader/internal/scrapper"
)

type Chain struct {
	Links        []Link
	ChainMessage *ChainMessage
}

type ChainMessage struct {
	ScrapType   scrapper.ScrapType
	URL         bandcamp.BandcampURL
	StorageType string
}

func NewChain(links []Link) *Chain {
	chainMessage := &ChainMessage{
		ScrapType: scrapper.Undefined,
	}
	lastElement := len(links) - 1
	for i, h := range links {
		if i < lastElement {
			nextHandler := links[i+1]
			h.SetNext(nextHandler)
		}
	}
	if links[0] != nil {
		links[0].Handle(chainMessage)
	}
	return &Chain{
		Links:        links,
		ChainMessage: chainMessage,
	}
}
