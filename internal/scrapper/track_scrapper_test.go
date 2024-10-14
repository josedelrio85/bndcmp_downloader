package scrapper

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"net/url"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/josedelrio85/bndcmp_downloader/internal/album_catalog"
	"github.com/josedelrio85/bndcmp_downloader/internal/bandcamp"
	"github.com/josedelrio85/bndcmp_downloader/internal/model"
	"github.com/stretchr/testify/suite"
	html "golang.org/x/net/html"
)

//go:embed resources/node_example.xml
var validExample string

//go:embed resources/invalid_node_example.xml
var invalidExample string

//go:embed resources/tralbum.json
var validJSONExample string

func TestTrackScrapper(t *testing.T) {
	suite.Run(t, new(TestTrackScrapperSuite))
}

type TestTrackScrapperSuite struct {
	suite.Suite
	controller      *gomock.Controller
	mockHttpClient  *MockRetriever
	mockParseClient *MockParser
	mockSaveClient  *MockSaver
	trackURL        *url.URL
	trackScrapper   *TrackScrapper
	albumCatalog    *album_catalog.MockAlbumCatalog
}

func (s *TestTrackScrapperSuite) SetupTest() {
	s.controller = gomock.NewController(s.T())
	s.mockHttpClient = NewMockRetriever(s.controller)
	s.mockParseClient = NewMockParser(s.controller)
	s.mockSaveClient = NewMockSaver(s.controller)
	trackURL, err := url.Parse("https://kinggizzard.bandcamp.com/track/elbow")
	s.NoError(err)
	s.trackURL = trackURL
	s.albumCatalog = album_catalog.NewMockAlbumCatalog(s.controller)
	s.trackScrapper = NewTrackScrapper(s.mockHttpClient, s.mockParseClient, s.mockSaveClient, s.albumCatalog)
}

func (s *TestTrackScrapperSuite) TearDownTest() {
	s.controller.Finish()
}

func (s *TestTrackScrapperSuite) TestRetrieve_Success() {
	mockResponse := []byte("mock response data")
	s.mockHttpClient.EXPECT().Retrieve(s.trackURL.String()).Return(bytes.NewReader(mockResponse), nil)

	reader, err := s.trackScrapper.Retrieve(s.trackURL.String())

	s.NoError(err)
	s.NotNil(reader)
}

func (s *TestTrackScrapperSuite) TestRetrieve_Error() {
	expectedError := errors.New("failed to retrieve track")
	s.mockHttpClient.EXPECT().Retrieve(s.trackURL.String()).Return(nil, expectedError)

	reader, err := s.trackScrapper.Retrieve(s.trackURL.String())

	s.Error(err)
	s.Equal(expectedError, err)
	s.Nil(reader)
}

func (s *TestTrackScrapperSuite) TestParse_Success() {
	mockResponse := []byte("mock response data")
	mockReader := bytes.NewReader(mockResponse)
	s.mockParseClient.EXPECT().Parse(mockReader).Return(&html.Node{}, nil)

	node, err := s.trackScrapper.Parse(mockReader)

	s.NoError(err)
	s.NotNil(node)
	s.Assert().IsType(&html.Node{}, node)
}

func (s *TestTrackScrapperSuite) TestParse_Error() {
	mockResponse := []byte("mock response data")
	mockReader := bytes.NewReader(mockResponse)
	expectedError := errors.New("failed to parse HTML")
	s.mockParseClient.EXPECT().Parse(mockReader).Return(nil, expectedError)

	node, err := s.trackScrapper.Parse(mockReader)

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

	err = s.trackScrapper.Find(nodes)
	s.NoError(err)
	s.Equal("Elbow", s.trackScrapper.Track.Title)
	s.Equal("https://kinggizzard.bandcamp.com/track/elbow", s.trackScrapper.Track.URL)
	s.Equal("https://t4.bcbits.com/stream/b77ce644d30f5a71778080be8c194c19/mp3-128/3749823254?p=0&ts=1728551843&t=dd8cc7cd9d747ac5be9c0a202fea450a5aa08944&token=1728551843_656b69850113f6ea23cd1e4321e6d148a256413b", s.trackScrapper.Track.DownloadURL)
}

func (s *TestTrackScrapperSuite) TestFind_Error() {
	dataAsBytes := []byte(invalidExample)
	mockReader := bytes.NewReader(dataAsBytes)
	nodes, err := html.Parse(mockReader)
	s.NoError(err)
	s.NotNil(nodes)

	err = s.trackScrapper.Find(nodes)
	s.Error(err)
	s.Equal(&model.Track{}, s.trackScrapper.Track)
	s.Contains(err.Error(), "invalid character")
}

func (s *TestTrackScrapperSuite) TestSave_Success() {
	mockReader := bytes.NewReader([]byte("mock response data"))
	track := &model.Track{}

	s.mockSaveClient.EXPECT().Save(mockReader, track).Return(nil)

	err := s.trackScrapper.Save(mockReader, track)
	s.NoError(err)
}

func (s *TestTrackScrapperSuite) TestSave_Error() {
	mockReader := bytes.NewReader([]byte("mock response data"))
	track := &model.Track{}

	mockedError := errors.New("failed to save file")
	s.mockSaveClient.EXPECT().Save(mockReader, track).Return(mockedError)

	err := s.trackScrapper.Save(mockReader, track)
	s.Error(err)
	s.Equal(mockedError, err)
}

func (s *TestTrackScrapperSuite) TestExecute_Success() {
	var trAlbum bandcamp.TrAlbum
	err := json.Unmarshal([]byte(validJSONExample), &trAlbum)
	if err != nil {
		s.T().Fatal(err)
	}
	downloadURL := trAlbum.Trackinfo[0].File.Mp3128

	mockReader := bytes.NewReader([]byte(validExample))
	s.mockHttpClient.EXPECT().Retrieve(s.trackURL.String()).Return(mockReader, nil)

	mockNode, _ := html.Parse(bytes.NewReader([]byte(validExample)))
	s.mockParseClient.EXPECT().Parse(mockReader).Return(mockNode, nil)

	expectedMapDir := make(map[string]bool)
	s.albumCatalog.EXPECT().GetMapDir().Return(&expectedMapDir).Times(2)

	mockMP3Reader := bytes.NewReader([]byte("mock mp3 data"))
	s.mockHttpClient.EXPECT().Retrieve(downloadURL).Return(mockMP3Reader, nil)
	s.trackScrapper.Track = trAlbum.ToTrack()
	s.mockSaveClient.EXPECT().Save(mockMP3Reader, s.trackScrapper.Track).Return(nil)
	s.albumCatalog.EXPECT().Update(s.trackScrapper.generateFilePath()).Return()

	err = s.trackScrapper.Execute(s.trackURL)

	s.NoError(err)
	s.Equal("Elbow", s.trackScrapper.Track.Title)
	s.Equal("https://kinggizzard.bandcamp.com/track/elbow", s.trackScrapper.Track.URL)
	s.Equal("https://t4.bcbits.com/stream/b77ce644d30f5a71778080be8c194c19/mp3-128/3749823254?p=0&ts=1728551843&t=dd8cc7cd9d747ac5be9c0a202fea450a5aa08944&token=1728551843_656b69850113f6ea23cd1e4321e6d148a256413b", s.trackScrapper.Track.DownloadURL)
}

func (s *TestTrackScrapperSuite) TestExecute_RetrieveError() {
	mockError := errors.New("retrieve error")
	s.mockHttpClient.EXPECT().Retrieve(s.trackURL.String()).Return(nil, mockError)

	err := s.trackScrapper.Execute(s.trackURL)

	s.Error(err)
	s.Equal(mockError, err)
}

func (s *TestTrackScrapperSuite) TestExecute_ParseError() {
	mockReader := bytes.NewReader([]byte(validExample))
	s.mockHttpClient.EXPECT().Retrieve(s.trackURL.String()).Return(mockReader, nil)

	mockError := errors.New("parse error")
	s.mockParseClient.EXPECT().Parse(mockReader).Return(nil, mockError)

	err := s.trackScrapper.Execute(s.trackURL)

	s.Error(err)
	s.Equal(mockError, err)
}

func (s *TestTrackScrapperSuite) TestExecute_FindError() {
	mockReader := bytes.NewReader([]byte(invalidExample))
	s.mockHttpClient.EXPECT().Retrieve(s.trackURL.String()).Return(mockReader, nil)

	mockNode, _ := html.Parse(mockReader)
	s.mockParseClient.EXPECT().Parse(mockReader).Return(mockNode, nil)

	err := s.trackScrapper.Execute(s.trackURL)

	s.Error(err)
	s.Contains(err.Error(), "invalid character")
}

func (s *TestTrackScrapperSuite) TestExecute_SaveError() {
	var trAlbum bandcamp.TrAlbum
	err := json.Unmarshal([]byte(validJSONExample), &trAlbum)
	if err != nil {
		s.T().Fatal(err)
	}
	downloadURL := trAlbum.Trackinfo[0].File.Mp3128
	s.trackScrapper.Track = trAlbum.ToTrack()

	mockReader := bytes.NewReader([]byte(validExample))
	s.mockHttpClient.EXPECT().Retrieve(s.trackURL.String()).Return(mockReader, nil)

	mockNode, _ := html.Parse(mockReader)
	s.mockParseClient.EXPECT().Parse(mockReader).Return(mockNode, nil)

	expectedMapDir := make(map[string]bool)
	s.albumCatalog.EXPECT().GetMapDir().Return(&expectedMapDir)

	mockMP3Reader := bytes.NewReader([]byte("mock mp3 data"))
	s.mockHttpClient.EXPECT().Retrieve(downloadURL).Return(mockMP3Reader, nil)

	mockError := errors.New("save error")
	s.mockSaveClient.EXPECT().Save(mockMP3Reader, s.trackScrapper.Track).Return(mockError)

	err = s.trackScrapper.Execute(s.trackURL)

	s.Error(err)
	s.Equal(mockError, err)
}

func (s *TestTrackScrapperSuite) TestFind_NoDataTralbum() {
	mockNode := &html.Node{
		Type: html.ElementNode,
		Data: "div",
	}

	err := s.trackScrapper.Find(mockNode)

	s.NoError(err)
	s.Equal(&model.Track{}, s.trackScrapper.Track)
}

func (s *TestTrackScrapperSuite) TestFind_InvalidJSON() {
	mockNode := &html.Node{
		Type: html.ElementNode,
		Data: "script",
		Attr: []html.Attribute{
			{
				Key: "data-tralbum",
				Val: "{invalid json}",
			},
		},
	}

	err := s.trackScrapper.Find(mockNode)

	s.Error(err)
	s.Contains(err.Error(), "invalid character")
}

func (s *TestTrackScrapperSuite) TestIsDownloaded_True() {
	filePath := "Artist/Album/01 - Track.mp3"
	expectedMapDir := map[string]bool{filePath: true}
	s.albumCatalog.EXPECT().GetMapDir().Return(&expectedMapDir)
	s.trackScrapper.Track = &model.Track{
		Artist:      "Artist",
		Album:       toPointer("Album"),
		Title:       "Track",
		TrackNumber: 1,
	}

	result := s.trackScrapper.isDownloaded()

	s.True(result)
}

func (s *TestTrackScrapperSuite) TestIsDownloaded_False() {
	expectedMapDir := make(map[string]bool)
	s.albumCatalog.EXPECT().GetMapDir().Return(&expectedMapDir)

	s.trackScrapper.Track = &model.Track{
		Artist:      "Artist",
		Album:       toPointer("Album"),
		Title:       "Track",
		TrackNumber: 1,
	}

	result := s.trackScrapper.isDownloaded()

	s.False(result)
}

func (s *TestTrackScrapperSuite) TestGenerateFilePath_WithAlbum() {
	s.trackScrapper.Track = &model.Track{
		Artist:      "Artist",
		Album:       toPointer("Album"),
		Title:       "Track",
		TrackNumber: 1,
	}

	result := s.trackScrapper.generateFilePath()

	s.Equal("Artist/Album/01 - Track.mp3", result)
}

func (s *TestTrackScrapperSuite) TestGenerateFilePath_WithoutAlbum() {
	s.trackScrapper.Track = &model.Track{
		Artist:      "Artist",
		Title:       "Track",
		TrackNumber: 1,
	}

	result := s.trackScrapper.generateFilePath()

	s.Equal("Artist/01 - Track.mp3", result)
}

/*
	func (s *TestTrackScrapperSuite) TestUpdateDownloadedTracks() {
		// Initialize the downloadedTracks map
		// s.trackScrapper.albumCatalog.GetMapDir() = &map[string]bool{}

		// Set up the Track
		s.trackScrapper.Track = &model.Track{
			Artist:      "Artist",
			Album:       toPointer("Album"),
			Title:       "Track",
			TrackNumber: 1,
		}

		// Call the method
		s.trackScrapper.updateDownloadedTracks()

		// Check if the track was added to the map
		expectedPath := "Artist/Album/01 - Track.mp3"
		_, exists := (*&s.trackScrapper.albumCatalog.GetMapDir())[expectedPath]
		s.True(exists, "The track should be marked as downloaded")

		// Verify the value is set to true
		s.True((*s.trackScrapper.downloadedTracks)[expectedPath], "The value for the track should be true")

		// Check the length of the map
		s.Equal(1, len(*s.trackScrapper.downloadedTracks), "The map should contain exactly one entry")
	}
*/
func toPointer(s string) *string {
	return &s
}
