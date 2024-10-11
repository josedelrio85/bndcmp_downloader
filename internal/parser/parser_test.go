package parser

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/net/html"
)

func toPointer(s string) *string {
	return &s
}

func TestParseClient(t *testing.T) {
	suite.Run(t, new(TestParseClientSuite))
}

type TestParseClientSuite struct {
	suite.Suite
	parseClient *ParseClient
}

func (s *TestParseClientSuite) SetupTest() {
	s.parseClient = NewParseClient()
}

func (s *TestParseClientSuite) TestParseClient_Parse() {
	tests := []struct {
		name          string
		input         *string
		expectedError bool
	}{
		{
			name:          "Valid HTML",
			input:         toPointer("<html><body><h1>Hello, World!</h1></body></html>"),
			expectedError: false,
		},
		{
			name:          "Empty input",
			input:         toPointer(""),
			expectedError: false,
		},
		{
			name:          "Invalid HTML",
			input:         nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.input != nil {
				reader := strings.NewReader(*tt.input)
				node, err := s.parseClient.Parse(reader)
				if tt.expectedError {
					s.Error(err)
					s.Nil(node)
				} else {
					s.NoError(err)
					s.NotNil(node)
					s.IsType(&html.Node{}, node)
				}
			} else {
				node, err := s.parseClient.Parse(&ErrorReader{Err: io.ErrUnexpectedEOF})
				s.Error(err)
				s.Nil(node)
			}

		})
	}
}

func (s *TestParseClientSuite) TestParseClient_Parse_ReaderError() {
	errorReader := &ErrorReader{Err: io.ErrUnexpectedEOF}
	node, err := s.parseClient.Parse(errorReader)

	s.Error(err)
	s.Nil(node)
	s.Equal(io.ErrUnexpectedEOF, err)
}

// ErrorReader is a custom io.Reader that always returns an error
type ErrorReader struct {
	Err error
}

func (er *ErrorReader) Read(p []byte) (n int, err error) {
	return 0, er.Err
}
