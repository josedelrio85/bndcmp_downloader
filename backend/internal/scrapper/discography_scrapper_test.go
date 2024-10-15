package scrapper

import (
	"bytes"
	_ "embed"
	"errors"
	"net/url"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/josedelrio85/bndcmp_downloader/internal/album_catalog"
	model "github.com/josedelrio85/bndcmp_downloader/internal/model"
	"github.com/stretchr/testify/suite"
	html "golang.org/x/net/html"
)

//go:embed resources/valid_discography_example.html
var validDiscographyExample string

func TestDiscographyScrapper(t *testing.T) {
	suite.Run(t, new(TestDiscographyScrapperSuite))
}

type mockAlbumScrapper struct {
	ExecuteFunc  func() error
	ExecuteCalls int
	URL          *url.URL
}

func (m *mockAlbumScrapper) Execute(url *url.URL) error {
	m.ExecuteCalls++
	return m.ExecuteFunc()
}

type TestDiscographyScrapperSuite struct {
	suite.Suite
	controller          *gomock.Controller
	mockHttpClient      *MockRetriever
	mockParseClient     *MockParser
	mockSaveClient      *MockSaver
	discographyURL      *url.URL
	DiscographyScrapper *DiscographyScrapper
	albumCatalog        *album_catalog.MockAlbumCatalog
}

func (s *TestDiscographyScrapperSuite) SetupTest() {
	s.controller = gomock.NewController(s.T())
	s.mockHttpClient = NewMockRetriever(s.controller)
	s.mockParseClient = NewMockParser(s.controller)
	s.mockSaveClient = NewMockSaver(s.controller)
	s.discographyURL = &url.URL{
		Scheme: "https",
		Host:   "kinggizzard.bandcamp.com",
		Path:   "/music",
	}
	s.albumCatalog = album_catalog.NewMockAlbumCatalog(s.controller)
	s.DiscographyScrapper = NewDiscographyScrapper(s.mockHttpClient, s.mockParseClient, s.mockSaveClient, s.albumCatalog)
}

func (s *TestDiscographyScrapperSuite) TearDownTest() {
	s.controller.Finish()
}

func (s *TestDiscographyScrapperSuite) TestRetrieve_Success() {
	mockResponse := []byte("mock response data")
	s.mockHttpClient.EXPECT().Retrieve(s.discographyURL.String()).Return(bytes.NewReader(mockResponse), nil)

	reader, err := s.DiscographyScrapper.Retrieve(s.discographyURL.String())

	s.NoError(err)
	s.NotNil(reader)
}

func (s *TestDiscographyScrapperSuite) TestRetrieve_Error() {
	expectedError := errors.New("failed to retrieve track")
	s.mockHttpClient.EXPECT().Retrieve(s.discographyURL.String()).Return(nil, expectedError)

	reader, err := s.DiscographyScrapper.Retrieve(s.discographyURL.String())

	s.Error(err)
	s.Equal(expectedError, err)
	s.Nil(reader)
}

func (s *TestDiscographyScrapperSuite) TestParse_Success() {
	mockResponse := []byte("mock response data")
	mockReader := bytes.NewReader(mockResponse)
	s.mockParseClient.EXPECT().Parse(mockReader).Return(&html.Node{}, nil)

	node, err := s.DiscographyScrapper.Parse(mockReader)

	s.NoError(err)
	s.NotNil(node)
	s.Assert().IsType(&html.Node{}, node)
}

func (s *TestDiscographyScrapperSuite) TestParse_Error() {
	mockResponse := []byte("mock response data")
	mockReader := bytes.NewReader(mockResponse)
	expectedError := errors.New("failed to parse HTML")
	s.mockParseClient.EXPECT().Parse(mockReader).Return(nil, expectedError)

	node, err := s.DiscographyScrapper.Parse(mockReader)

	s.Error(err)
	s.Equal(expectedError, err)
	s.Nil(node)
}

func (s *TestDiscographyScrapperSuite) Test_find_Success() {
	dataAsBytes := []byte(validDiscographyExample)
	mockReader := bytes.NewReader(dataAsBytes)
	nodes, err := html.Parse(mockReader)
	s.NoError(err)
	s.NotNil(nodes)

	expectedAlbumList := []string{
		"/album/live-at-levitation-16",
		"/album/live-in-milwaukee-19",
		"/album/butterfly-3000",
		"/album/live-in-sydney-21",
		"/album/live-in-melbourne-21",
		"/album/l-w",
		"/album/live-in-london-19",
		"/album/teenage-gizzard",
		"/album/k-g",
		"https://kinggizzard.bandcamp.com/album/live-in-san-francisco-16",
		"/album/live-in-asheville-19",
		"/album/demos-vol-1-vol-2",
		"https://kinggizzard.bandcamp.com/album/chunky-shrapnel",
		"/album/live-in-brussels-19",
		"/album/live-in-adelaide-19",
		"/album/live-in-paris-19",
		"https://kinggizzard.bandcamp.com/album/infest-the-rats-nest-2",
		"/album/fishing-for-fishies",
		"/album/gumboot-soup",
		"/album/polygondwanaland",
		"/album/sketches-of-brunswick-east",
		"/album/murder-of-the-universe",
		"/album/flying-microtonal-banana",
		"/album/nonagon-infinity-2",
		"/album/paper-m-ch-dream-balloon",
		"/album/quarters",
		"/album/im-in-your-mind-fuzz",
		"/album/oddments",
		"/album/float-along-fill-your-lungs",
		"/album/eyes-like-the-sky",
		"/album/12-bar-bruise",
		"/album/willoughbys-beach-ep",
	}

	err = s.DiscographyScrapper.find(nodes)
	s.NoError(err)
	s.Equal(expectedAlbumList, s.DiscographyScrapper.AlbumList)
}

func (s *TestDiscographyScrapperSuite) Test_processAlbumList_Success() {
	dataAsBytes := []byte(validDiscographyExample)
	mockReader := bytes.NewReader(dataAsBytes)
	nodes, err := html.Parse(mockReader)
	s.NoError(err)
	s.NotNil(nodes)
	err = s.DiscographyScrapper.find(nodes)
	s.NoError(err)

	expectedDiscographyList := []string{
		"/album/live-at-levitation-16", "/album/live-in-milwaukee-19", "/album/butterfly-3000", "/album/live-in-sydney-21", "/album/live-in-melbourne-21", "/album/l-w", "/album/live-in-london-19", "/album/teenage-gizzard", "/album/k-g", "/album/live-in-asheville-19", "/album/demos-vol-1-vol-2", "/album/live-in-brussels-19", "/album/live-in-adelaide-19", "/album/live-in-paris-19", "/album/fishing-for-fishies", "/album/gumboot-soup", "/album/polygondwanaland", "/album/sketches-of-brunswick-east", "/album/murder-of-the-universe", "/album/flying-microtonal-banana", "/album/nonagon-infinity-2", "/album/paper-m-ch-dream-balloon", "/album/quarters", "/album/im-in-your-mind-fuzz", "/album/oddments", "/album/float-along-fill-your-lungs", "/album/eyes-like-the-sky", "/album/12-bar-bruise", "/album/willoughbys-beach-ep",
	}

	output := s.DiscographyScrapper.processAlbumList()
	s.Equal(expectedDiscographyList, output)
}

func (s *TestDiscographyScrapperSuite) Test_processAlbumList_NoResults() {
	dataAsBytes := []byte(invalidExample)
	mockReader := bytes.NewReader(dataAsBytes)
	nodes, err := html.Parse(mockReader)
	s.NoError(err)
	s.NotNil(nodes)
	err = s.DiscographyScrapper.find(nodes)
	s.NoError(err)

	output := s.DiscographyScrapper.processAlbumList()
	s.Equal(0, len(output))
}

func (s *TestDiscographyScrapperSuite) TestFind_Success() {
	dataAsBytes := []byte(validDiscographyExample)
	mockReader := bytes.NewReader(dataAsBytes)
	nodes, err := html.Parse(mockReader)
	s.NoError(err)
	s.NotNil(nodes)

	expectedDiscographyList := []string{
		"/album/live-at-levitation-16", "/album/live-in-milwaukee-19", "/album/butterfly-3000", "/album/live-in-sydney-21", "/album/live-in-melbourne-21", "/album/l-w", "/album/live-in-london-19", "/album/teenage-gizzard", "/album/k-g", "/album/live-in-asheville-19", "/album/demos-vol-1-vol-2", "/album/live-in-brussels-19", "/album/live-in-adelaide-19", "/album/live-in-paris-19", "/album/fishing-for-fishies", "/album/gumboot-soup", "/album/polygondwanaland", "/album/sketches-of-brunswick-east", "/album/murder-of-the-universe", "/album/flying-microtonal-banana", "/album/nonagon-infinity-2", "/album/paper-m-ch-dream-balloon", "/album/quarters", "/album/im-in-your-mind-fuzz", "/album/oddments", "/album/float-along-fill-your-lungs", "/album/eyes-like-the-sky", "/album/12-bar-bruise", "/album/willoughbys-beach-ep",
	}

	err = s.DiscographyScrapper.Find(nodes)
	s.NoError(err)
	s.Equal(expectedDiscographyList, s.DiscographyScrapper.AlbumList)
}

func (s *TestDiscographyScrapperSuite) TestFind_Error() {
	dataAsBytes := []byte(invalidExample)
	mockReader := bytes.NewReader(dataAsBytes)
	nodes, err := html.Parse(mockReader)
	s.NoError(err)
	s.NotNil(nodes)

	err = s.DiscographyScrapper.Find(nodes)
	s.NoError(err)
	s.Equal(0, len(s.DiscographyScrapper.AlbumList))
}

func (s *TestDiscographyScrapperSuite) TestSave() {
	mockReader := bytes.NewReader([]byte("mock response data"))

	err := s.DiscographyScrapper.Save(mockReader, &model.Track{})
	s.NoError(err)
}

func (s *TestDiscographyScrapperSuite) TestExecute_Success() {
	mockExecuteClient := &mockAlbumScrapper{
		ExecuteFunc: func() error {
			return nil
		},
	}
	s.DiscographyScrapper.executeClient = func(httpClient Retriever, parseClient Parser, saveClient Saver, albumCatalog album_catalog.AlbumCatalog) Executer {
		return mockExecuteClient
	}

	mockReader := bytes.NewReader([]byte(validDiscographyExample))
	s.mockHttpClient.EXPECT().Retrieve(s.discographyURL.String()).Return(mockReader, nil)

	mockNode, _ := html.Parse(bytes.NewReader([]byte(validDiscographyExample)))
	s.mockParseClient.EXPECT().Parse(mockReader).Return(mockNode, nil)

	err := s.DiscographyScrapper.Execute(s.discographyURL)

	s.NoError(err)
	s.Equal(len(s.DiscographyScrapper.AlbumList), mockExecuteClient.ExecuteCalls)
}

func (s *TestDiscographyScrapperSuite) TestExecute_RetrieveError() {
	mockError := errors.New("retrieve error")
	s.mockHttpClient.EXPECT().Retrieve(s.discographyURL.String()).Return(nil, mockError)

	err := s.DiscographyScrapper.Execute(s.discographyURL)

	s.Error(err)
	s.Equal(mockError, err)
}

func (s *TestDiscographyScrapperSuite) TestExecute_ParseError() {
	mockReader := bytes.NewReader([]byte(validExample))
	s.mockHttpClient.EXPECT().Retrieve(s.discographyURL.String()).Return(mockReader, nil)

	mockError := errors.New("parse error")
	s.mockParseClient.EXPECT().Parse(mockReader).Return(nil, mockError)

	err := s.DiscographyScrapper.Execute(s.discographyURL)

	s.Error(err)
	s.Equal(mockError, err)
}

func (s *TestDiscographyScrapperSuite) TestExecute_FindError() {
	mockReader := bytes.NewReader([]byte(invalidExample))
	s.mockHttpClient.EXPECT().Retrieve(s.discographyURL.String()).Return(mockReader, nil)

	mockNode, _ := html.Parse(mockReader)
	s.mockParseClient.EXPECT().Parse(mockReader).Return(mockNode, nil)

	err := s.DiscographyScrapper.Execute(s.discographyURL)

	s.NoError(err)
	s.Equal(0, len(s.DiscographyScrapper.AlbumList))
}

func (s *TestDiscographyScrapperSuite) TestExecute_SaveError() {
	mockedError := errors.New("album scrapper error")

	mockExecuteClient := &mockAlbumScrapper{
		ExecuteFunc: func() error {
			return mockedError
		},
	}
	s.DiscographyScrapper.executeClient = func(httpClient Retriever, parseClient Parser, saveClient Saver, albumCatalog album_catalog.AlbumCatalog) Executer {
		return mockExecuteClient
	}

	mockReader := bytes.NewReader([]byte(validDiscographyExample))
	s.mockHttpClient.EXPECT().Retrieve(s.discographyURL.String()).Return(mockReader, nil)

	mockNode, _ := html.Parse(mockReader)
	s.mockParseClient.EXPECT().Parse(mockReader).Return(mockNode, nil)

	err := s.DiscographyScrapper.Execute(s.discographyURL)

	s.Error(err)
	s.Equal(mockedError, err)
}
