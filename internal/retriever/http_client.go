package retriever

import (
	"io"
	"net/http"
)

type httpClient struct {
}

func NewHttpClient() *httpClient {
	return &httpClient{}
}

func (h *httpClient) Retrieve(url string) (io.Reader, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}
