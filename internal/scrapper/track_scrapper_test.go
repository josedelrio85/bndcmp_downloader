package scrapper

import (
	"bytes"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestTrackScrapper(t *testing.T) {
	suite.Run(t, new(TestTrackScrapperSuite))
}

type TestTrackScrapperSuite struct {
	suite.Suite
	controller     *gomock.Controller
	mockHttpClient *MockRetriever
	trackURL       string
	scrapper       *TrackScrapper
}

func (s *TestTrackScrapperSuite) SetupTest() {
	s.controller = gomock.NewController(s.T())
	s.mockHttpClient = NewMockRetriever(s.controller)
	s.trackURL = "https://kinggizzard.bandcamp.com/track/elbow"
	s.scrapper = NewTrackScrapper(s.trackURL, s.mockHttpClient)
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
