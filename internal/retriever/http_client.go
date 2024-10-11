package retriever

import (
	"io"
	"net/http"
)

type HttpClient struct {
}

func NewHttpClient() *HttpClient {
	return &HttpClient{}
}

func (h *HttpClient) Retrieve(url string) (io.Reader, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}
