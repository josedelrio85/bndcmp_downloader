package scrapper

import (
	"bytes"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	html "golang.org/x/net/html"
)

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

	s.Require().NoError(err)
	s.Require().NotNil(reader)
}

func (s *TestTrackScrapperSuite) TestRetrieve_Error() {
	expectedError := errors.New("failed to retrieve track")
	s.mockHttpClient.EXPECT().Retrieve(s.trackURL).Return(nil, expectedError)

	reader, err := s.scrapper.Retrieve(s.trackURL)

	s.Require().Error(err)
	s.Require().Equal(expectedError, err)
	s.Require().Nil(reader)
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
