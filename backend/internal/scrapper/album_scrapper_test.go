package scrapper

import (
	"bytes"
	_ "embed"
	"errors"
	"net/url"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/josedelrio85/bndcmp_downloader/internal/album_catalog"
	"github.com/josedelrio85/bndcmp_downloader/internal/model"
	"github.com/stretchr/testify/suite"
	html "golang.org/x/net/html"
)

//go:embed resources/valid_album_example.html
var validAlbumExample string

func TestAlbumScrapper(t *testing.T) {
	suite.Run(t, new(TestalbumScrapperSuite))
}

type mockTrackScrapper struct {
	ExecuteFunc  func() error
	ExecuteCalls int
	URL          string
}

func (m *mockTrackScrapper) Execute(url *url.URL) error {
	m.ExecuteCalls++
	return m.ExecuteFunc()
}

type TestalbumScrapperSuite struct {
	suite.Suite
	controller      *gomock.Controller
	mockHttpClient  *MockRetriever
	mockParseClient *MockParser
	mockSaveClient  *MockSaver
	albumURL        *url.URL
	albumScrapper   *AlbumScrapper
	albumCatalog    *album_catalog.MockAlbumCatalog
}

func (s *TestalbumScrapperSuite) SetupTest() {
	s.controller = gomock.NewController(s.T())
	s.mockHttpClient = NewMockRetriever(s.controller)
	s.mockParseClient = NewMockParser(s.controller)
	s.mockSaveClient = NewMockSaver(s.controller)
	s.albumURL = &url.URL{
		Scheme: "https",
		Host:   "kinggizzard.bandcamp.com",
		Path:   "/album/12-bar-bruise",
	}
	s.albumCatalog = album_catalog.NewMockAlbumCatalog(s.controller)
	s.albumScrapper = NewAlbumScrapper(s.mockHttpClient, s.mockParseClient, s.mockSaveClient, s.albumCatalog)
}

func (s *TestalbumScrapperSuite) TearDownTest() {
	s.controller.Finish()
}

func (s *TestalbumScrapperSuite) TestRetrieve_Success() {
	mockResponse := []byte("mock response data")
	s.mockHttpClient.EXPECT().Retrieve(s.albumURL.String()).Return(bytes.NewReader(mockResponse), nil)

	reader, err := s.albumScrapper.Retrieve(s.albumURL.String())

	s.NoError(err)
	s.NotNil(reader)
}

func (s *TestalbumScrapperSuite) TestRetrieve_Error() {
	expectedError := errors.New("failed to retrieve track")
	s.mockHttpClient.EXPECT().Retrieve(s.albumURL.String()).Return(nil, expectedError)

	reader, err := s.albumScrapper.Retrieve(s.albumURL.String())

	s.Error(err)
	s.Equal(expectedError, err)
	s.Nil(reader)
}

func (s *TestalbumScrapperSuite) TestParse_Success() {
	mockResponse := []byte("mock response data")
	mockReader := bytes.NewReader(mockResponse)
	s.mockParseClient.EXPECT().Parse(mockReader).Return(&html.Node{}, nil)

	node, err := s.albumScrapper.Parse(mockReader)

	s.NoError(err)
	s.NotNil(node)
	s.Assert().IsType(&html.Node{}, node)
}

func (s *TestalbumScrapperSuite) TestParse_Error() {
	mockResponse := []byte("mock response data")
	mockReader := bytes.NewReader(mockResponse)
	expectedError := errors.New("failed to parse HTML")
	s.mockParseClient.EXPECT().Parse(mockReader).Return(nil, expectedError)

	node, err := s.albumScrapper.Parse(mockReader)

	s.Error(err)
	s.Equal(expectedError, err)
	s.Nil(node)
}

func (s *TestalbumScrapperSuite) Test_find_Success() {
	dataAsBytes := []byte(validAlbumExample)
	mockReader := bytes.NewReader(dataAsBytes)
	nodes, err := html.Parse(mockReader)
	s.NoError(err)
	s.NotNil(nodes)

	expectedTrackList := []string{
		"/track/elbow",
		"/track/elbow",
		"/track/elbow?action=download",
		"/track/muckraker",
		"/track/muckraker",
		"/track/muckraker?action=download",
		"/track/nein",
		"/track/nein",
		"/track/nein?action=download",
		"/track/12-bar-bruise",
		"/track/12-bar-bruise",
		"/track/12-bar-bruise?action=download",
		"/track/garage-liddiard",
		"/track/garage-liddiard",
		"/track/garage-liddiard?action=download",
		"/track/sam-cherrys-last-shot-2",
		"/track/sam-cherrys-last-shot-2",
		"/track/sam-cherrys-last-shot-2?action=download",
		"/track/high-hopes-low",
		"/track/high-hopes-low",
		"/track/high-hopes-low?action=download",
		"/track/cut-throat-boogie",
		"/track/cut-throat-boogie",
		"/track/cut-throat-boogie?action=download",
		"/track/bloody-ripper-3",
		"/track/bloody-ripper-3",
		"/track/bloody-ripper-3?action=download",
		"/track/uh-oh-i-called-mum",
		"/track/uh-oh-i-called-mum",
		"/track/uh-oh-i-called-mum?action=download",
		"/track/sea-of-trees",
		"/track/sea-of-trees",
		"/track/sea-of-trees?action=download",
		"/track/footy-footy",
		"/track/footy-footy",
		"/track/footy-footy?action=download",
	}

	err = s.albumScrapper.find(nodes)
	s.NoError(err)
	s.Equal(expectedTrackList, s.albumScrapper.TrackList)
}

func (s *TestalbumScrapperSuite) Test_processTrackList_Success() {
	dataAsBytes := []byte(validAlbumExample)
	mockReader := bytes.NewReader(dataAsBytes)
	nodes, err := html.Parse(mockReader)
	s.NoError(err)
	s.NotNil(nodes)
	err = s.albumScrapper.find(nodes)
	s.NoError(err)

	expectedTrackList := []string{
		"/track/elbow",
		"/track/muckraker",
		"/track/nein",
		"/track/12-bar-bruise",
		"/track/garage-liddiard",
		"/track/sam-cherrys-last-shot-2",
		"/track/high-hopes-low",
		"/track/cut-throat-boogie",
		"/track/bloody-ripper-3",
		"/track/uh-oh-i-called-mum",
		"/track/sea-of-trees",
		"/track/footy-footy",
	}

	output := s.albumScrapper.processTrackList()
	s.Equal(expectedTrackList, output)
}

func (s *TestalbumScrapperSuite) Test_processTrackList_NoResults() {
	dataAsBytes := []byte(invalidExample)
	mockReader := bytes.NewReader(dataAsBytes)
	nodes, err := html.Parse(mockReader)
	s.NoError(err)
	s.NotNil(nodes)
	err = s.albumScrapper.find(nodes)
	s.NoError(err)

	output := s.albumScrapper.processTrackList()
	s.Equal(0, len(output))
}

func (s *TestalbumScrapperSuite) TestFind_Success() {
	dataAsBytes := []byte(validAlbumExample)
	mockReader := bytes.NewReader(dataAsBytes)
	nodes, err := html.Parse(mockReader)
	s.NoError(err)
	s.NotNil(nodes)

	expectedTrackList := []string{
		"/track/elbow",
		"/track/muckraker",
		"/track/nein",
		"/track/12-bar-bruise",
		"/track/garage-liddiard",
		"/track/sam-cherrys-last-shot-2",
		"/track/high-hopes-low",
		"/track/cut-throat-boogie",
		"/track/bloody-ripper-3",
		"/track/uh-oh-i-called-mum",
		"/track/sea-of-trees",
		"/track/footy-footy",
	}

	err = s.albumScrapper.Find(nodes)
	s.NoError(err)
	s.Equal(expectedTrackList, s.albumScrapper.TrackList)
}

func (s *TestalbumScrapperSuite) TestFind_Error() {
	dataAsBytes := []byte(invalidExample)
	mockReader := bytes.NewReader(dataAsBytes)
	nodes, err := html.Parse(mockReader)
	s.NoError(err)
	s.NotNil(nodes)

	err = s.albumScrapper.Find(nodes)
	s.NoError(err)
	s.Equal(0, len(s.albumScrapper.TrackList))
}

func (s *TestalbumScrapperSuite) TestSave() {
	mockReader := bytes.NewReader([]byte("mock response data"))

	err := s.albumScrapper.Save(mockReader, &model.Track{})
	s.NoError(err)
}

func (s *TestalbumScrapperSuite) TestExecute_Success() {
	mockExecuteClient := &mockTrackScrapper{
		ExecuteFunc: func() error {
			return nil
		},
	}
	s.albumScrapper.executeClient = func(httpClient Retriever, parseClient Parser, saveClient Saver, albumCatalog album_catalog.AlbumCatalog) Executer {
		return mockExecuteClient
	}

	mockReader := bytes.NewReader([]byte(validAlbumExample))
	s.mockHttpClient.EXPECT().Retrieve(s.albumURL.String()).Return(mockReader, nil)

	mockNode, _ := html.Parse(bytes.NewReader([]byte(validAlbumExample)))
	s.mockParseClient.EXPECT().Parse(mockReader).Return(mockNode, nil)

	err := s.albumScrapper.Execute(s.albumURL)

	s.NoError(err)
	s.Equal(len(s.albumScrapper.TrackList), mockExecuteClient.ExecuteCalls)
}

func (s *TestalbumScrapperSuite) TestExecute_RetrieveError() {
	mockError := errors.New("retrieve error")
	s.mockHttpClient.EXPECT().Retrieve(s.albumURL.String()).Return(nil, mockError)

	err := s.albumScrapper.Execute(s.albumURL)

	s.Error(err)
	s.Equal(mockError, err)
}

func (s *TestalbumScrapperSuite) TestExecute_ParseError() {
	mockReader := bytes.NewReader([]byte(validExample))
	s.mockHttpClient.EXPECT().Retrieve(s.albumURL.String()).Return(mockReader, nil)

	mockError := errors.New("parse error")
	s.mockParseClient.EXPECT().Parse(mockReader).Return(nil, mockError)

	err := s.albumScrapper.Execute(s.albumURL)

	s.Error(err)
	s.Equal(mockError, err)
}

func (s *TestalbumScrapperSuite) TestExecute_FindError() {
	mockReader := bytes.NewReader([]byte(invalidExample))
	s.mockHttpClient.EXPECT().Retrieve(s.albumURL.String()).Return(mockReader, nil)

	mockNode, _ := html.Parse(mockReader)
	s.mockParseClient.EXPECT().Parse(mockReader).Return(mockNode, nil)

	err := s.albumScrapper.Execute(s.albumURL)

	s.NoError(err)
	s.Equal(0, len(s.albumScrapper.TrackList))
}

func (s *TestalbumScrapperSuite) TestExecute_SaveError() {
	mockedError := errors.New("track scrapper error")

	mockExecuteClient := &mockTrackScrapper{
		ExecuteFunc: func() error {
			return mockedError
		},
	}
	s.albumScrapper.executeClient = func(httpClient Retriever, parseClient Parser, saveClient Saver, albumCatalog album_catalog.AlbumCatalog) Executer {
		return mockExecuteClient
	}

	mockReader := bytes.NewReader([]byte(validAlbumExample))
	s.mockHttpClient.EXPECT().Retrieve(s.albumURL.String()).Return(mockReader, nil)

	mockNode, _ := html.Parse(mockReader)
	s.mockParseClient.EXPECT().Parse(mockReader).Return(mockNode, nil)

	err := s.albumScrapper.Execute(s.albumURL)

	s.Error(err)
	s.Equal(mockedError, err)
}
