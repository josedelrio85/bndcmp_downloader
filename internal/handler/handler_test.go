package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/josedelrio85/bndcmp_downloader/internal/scrapper"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	ctrl                    *gomock.Controller
	handler                 *HttpHandler
	mockDiscographyScrapper *scrapper.MockScrapper
	mockAlbumScrapper       *scrapper.MockScrapper
	mockTrackScrapper       *scrapper.MockScrapper
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (s *HandlerTestSuite) SetupTest() {
	baseFolder := "test_downloads"
	s.ctrl = gomock.NewController(s.T())

	s.mockDiscographyScrapper = scrapper.NewMockScrapper(s.ctrl)
	s.mockAlbumScrapper = scrapper.NewMockScrapper(s.ctrl)
	s.mockTrackScrapper = scrapper.NewMockScrapper(s.ctrl)

	s.handler = NewHttpHandler(
		baseFolder,
		s.mockDiscographyScrapper,
		s.mockAlbumScrapper,
		s.mockTrackScrapper,
	)
}

func (s *HandlerTestSuite) TestHealth() {
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	s.Require().NoError(err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.Health)

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	s.Equal("OK", rr.Body.String())
}

func (s *HandlerTestSuite) TestGetDiscography() {
	req, err := http.NewRequest("GET", "/api/v1/testartist", nil)
	s.Require().NoError(err)
	// Use the router to create a new request
	req = mux.SetURLVars(req, map[string]string{"artist": "testartist"})

	s.mockDiscographyScrapper.EXPECT().Execute(gomock.Any()).Return(nil).AnyTimes()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.GetDiscography)

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	s.Equal("Discography of testartist downloaded successfully", rr.Body.String())
}

func (s *HandlerTestSuite) TestGetAlbum() {
	req, err := http.NewRequest("GET", "/testartist/testalbum", nil)
	s.Require().NoError(err)
	req = mux.SetURLVars(req, map[string]string{"artist": "testartist", "album": "testalbum"})

	s.mockAlbumScrapper.EXPECT().Execute(gomock.Any()).Return(nil).AnyTimes()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.GetAlbum)

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	s.Equal("Album testalbum of testartist downloaded successfully", rr.Body.String())
}

func (s *HandlerTestSuite) TestGetTrack() {
	req, err := http.NewRequest("GET", "/testartist/track/testtrack", nil)
	s.Require().NoError(err)
	req = mux.SetURLVars(req, map[string]string{"artist": "testartist", "track": "testtrack"})

	s.mockTrackScrapper.EXPECT().Execute(gomock.Any()).Return(nil).AnyTimes()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.GetTrack)

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusOK, rr.Code)
	s.Equal("Track testtrack of testartist downloaded successfully", rr.Body.String())
}

func (s *HandlerTestSuite) TestGetDiscography_BadRequest() {
	req, err := http.NewRequest("GET", "/api/v1/", nil)
	s.Require().NoError(err)
	// Use the router to create a new request without artist
	req = mux.SetURLVars(req, map[string]string{})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.GetDiscography)

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
	s.Equal("Artist is required\n", rr.Body.String())
}

func (s *HandlerTestSuite) TestGetAlbum_BadRequest_NoArtist() {
	req, err := http.NewRequest("GET", "/testalbum", nil)
	s.Require().NoError(err)
	req = mux.SetURLVars(req, map[string]string{"album": "testalbum"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.GetAlbum)

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
	s.Equal("Artist is required\n", rr.Body.String())
}

func (s *HandlerTestSuite) TestGetAlbum_BadRequest_NoAlbum() {
	req, err := http.NewRequest("GET", "/testartist", nil)
	s.Require().NoError(err)
	req = mux.SetURLVars(req, map[string]string{"artist": "testartist"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.GetAlbum)

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
	s.Equal("Album is required\n", rr.Body.String())
}

func (s *HandlerTestSuite) TestGetTrack_BadRequest_NoArtist() {
	req, err := http.NewRequest("GET", "/track/testtrack", nil)
	s.Require().NoError(err)
	req = mux.SetURLVars(req, map[string]string{"track": "testtrack"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.GetTrack)

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
	s.Equal("Artist is required\n", rr.Body.String())
}

func (s *HandlerTestSuite) TestGetTrack_BadRequest_NoTrack() {
	req, err := http.NewRequest("GET", "/testartist/track", nil)
	s.Require().NoError(err)
	req = mux.SetURLVars(req, map[string]string{"artist": "testartist"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.handler.GetTrack)

	handler.ServeHTTP(rr, req)

	s.Equal(http.StatusBadRequest, rr.Code)
	s.Equal("Track is required\n", rr.Body.String())
}
