package retriever

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestHttpClient(t *testing.T) {
	suite.Run(t, new(TestHttpClientSuite))
}

type TestHttpClientSuite struct {
	suite.Suite
}

func (hc *TestHttpClientSuite) TestHttpClient_Retrieve() {
	tests := []struct {
		name           string
		serverResponse string
		serverStatus   int
		expectedError  bool
	}{
		{
			name:           "Successful retrieval",
			serverResponse: "Hello, World!",
			serverStatus:   http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "Server error",
			serverResponse: "Internal Server Error",
			serverStatus:   http.StatusInternalServerError,
			expectedError:  false, // The method doesn't check for non-200 status codes
		},
		{
			name:           "Empty response",
			serverResponse: "",
			serverStatus:   http.StatusOK,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		hc.Run(tt.name, func() {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// Create a new httpClient
			client := NewHttpClient()

			// Call the Retrieve method
			reader, err := client.Retrieve(server.URL)

			if tt.expectedError {
				hc.Error(err)
			}

			hc.NoError(err)
			if err == nil {
				hc.NotNil(reader)

				body, err := io.ReadAll(reader)
				hc.NoError(err)

				hc.Equal(tt.serverResponse, string(body))
			}
		})
	}
}

func (hc *TestHttpClientSuite) TestHttpClient_Retrieve_InvalidURL() {
	client := NewHttpClient()

	_, err := client.Retrieve("://invalid-url")

	hc.Error(err)
	hc.Contains(err.Error(), "invalid-url")
}

func (hc *TestHttpClientSuite) TestHttpClient_Retrieve_NetworkError() {
	client := NewHttpClient()

	_, err := client.Retrieve("http://non-existent-server.com")

	hc.Error(err)
}
