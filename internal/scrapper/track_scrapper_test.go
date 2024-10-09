package scrapper

import (
	"bytes"
	_ "embed"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/josedelrio85/bndcmp_downloader/internal/model"
	"github.com/stretchr/testify/suite"
	html "golang.org/x/net/html"
)

//go:embed node_example.xml
var validExample string

//go:embed invalid_node_example.xml
var invalidExample string

func TestTrackScrapper(t *testing.T) {
	suite.Run(t, new(TestTrackScrapperSuite))
}

type TestTrackScrapperSuite struct {
	suite.Suite
	controller      *gomock.Controller
	mockHttpClient  *MockRetriever
	mockParseClient *MockParser
	trackURL        string
	scrapper        *TrackScrapper
}

func (s *TestTrackScrapperSuite) SetupTest() {
	s.controller = gomock.NewController(s.T())
	s.mockHttpClient = NewMockRetriever(s.controller)
	s.mockParseClient = NewMockParser(s.controller)
	s.trackURL = "https://kinggizzard.bandcamp.com/track/elbow"
	s.scrapper = NewTrackScrapper(s.trackURL, s.mockHttpClient, s.mockParseClient)
}

func (s *TestTrackScrapperSuite) TearDownTest() {
	s.controller.Finish()
}

func (s *TestTrackScrapperSuite) TestRetrieve_Success() {
	mockResponse := []byte("mock response data")
	s.mockHttpClient.EXPECT().Retrieve(s.trackURL).Return(bytes.NewReader(mockResponse), nil)

	reader, err := s.scrapper.Retrieve(s.trackURL)

	s.NoError(err)
	s.NotNil(reader)
}

func (s *TestTrackScrapperSuite) TestRetrieve_Error() {
	expectedError := errors.New("failed to retrieve track")
	s.mockHttpClient.EXPECT().Retrieve(s.trackURL).Return(nil, expectedError)

	reader, err := s.scrapper.Retrieve(s.trackURL)

	s.Error(err)
	s.Equal(expectedError, err)
	s.Nil(reader)
}

func (s *TestTrackScrapperSuite) TestParse_Success() {
	mockResponse := []byte("mock response data")
	mockReader := bytes.NewReader(mockResponse)
	s.mockParseClient.EXPECT().Parse(mockReader).Return(&html.Node{}, nil)

	node, err := s.scrapper.Parse(mockReader)

	s.NoError(err)
	s.NotNil(node)
	s.Assert().IsType(&html.Node{}, node)
}

func (s *TestTrackScrapperSuite) TestParse_Error() {
	mockResponse := []byte("mock response data")
	mockReader := bytes.NewReader(mockResponse)
	expectedError := errors.New("failed to parse HTML")
	s.mockParseClient.EXPECT().Parse(mockReader).Return(nil, expectedError)

	node, err := s.scrapper.Parse(mockReader)

	s.Error(err)
	s.Equal(expectedError, err)
	s.Nil(node)
}

func (s *TestTrackScrapperSuite) TestFind_Success() {
	dataAsBytes := []byte(validExample)
	mockReader := bytes.NewReader(dataAsBytes)
	nodes, err := html.Parse(mockReader)
	s.NoError(err)
	s.NotNil(nodes)

	err = s.scrapper.Find(nodes)
	s.NoError(err)
	s.Equal("Elbow", s.scrapper.Track.Title)
	s.Equal("https://kinggizzard.bandcamp.com/track/elbow", s.scrapper.Track.URL)
	s.Equal("https://t4.bcbits.com/stream/b77ce644d30f5a71778080be8c194c19/mp3-128/3749823254?p=0&ts=1728551843&t=dd8cc7cd9d747ac5be9c0a202fea450a5aa08944&token=1728551843_656b69850113f6ea23cd1e4321e6d148a256413b", s.scrapper.Track.DownloadURL)
}

func (s *TestTrackScrapperSuite) TestFind_Error() {
	dataAsBytes := []byte(invalidExample)
	mockReader := bytes.NewReader(dataAsBytes)
	nodes, err := html.Parse(mockReader)
	s.NoError(err)
	s.NotNil(nodes)

	err = s.scrapper.Find(nodes)
	s.Error(err)
	s.Equal(&model.Track{}, s.scrapper.Track)
	s.Contains(err.Error(), "invalid character")
}
