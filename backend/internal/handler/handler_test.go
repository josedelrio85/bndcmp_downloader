package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

func (s *HandlerTestSuite) Test_getScrapper() {

	testCases := []struct {
		desc           string
		url            string
		expectedResult bool
	}{
		{
			desc:           "Valid discography URL",
			url:            "https://testartist.bandcamp.com/music",
			expectedResult: true,
		},
		{
			desc:           "Valid album URL",
			url:            "https://testartist.bandcamp.com/album/testalbum",
			expectedResult: true,
		},
		{
			desc:           "Valid track URL",
			url:            "https://testartist.bandcamp.com/track/testtrack",
			expectedResult: true,
		},
		{
			desc:           "Invalid URL",
			url:            "https://testartist.bandcamp.com/",
			expectedResult: false,
		},
		{
			desc:           "Invalid URL",
			url:            "https://testartist.bandcamp.com",
			expectedResult: false,
		},
	}

	for _, tt := range testCases {
		fmt.Println(tt.desc)

		scrapURL, err := url.Parse(tt.url)
		s.NoError(err)

		scrapper, err := s.handler.getScrapper(scrapURL)

		if tt.expectedResult {
			s.NoError(err)
			s.NotNil(scrapper)
		} else {
			s.Error(err)
			s.Nil(scrapper)
		}
	}
}

func (s *HandlerTestSuite) Test_isValidBandcampURL() {
	testCases := []struct {
		desc     string
		url      string
		expected bool
	}{
		{
			desc:     "Valid discography URL",
			url:      "https://testartist.bandcamp.com/music",
			expected: true,
		},
		{
			desc:     "Valid album URL",
			url:      "https://testartist.bandcamp.com/album/testalbum",
			expected: true,
		},
		{
			desc:     "Valid track URL",
			url:      "https://testartist.bandcamp.com/track/testtrack",
			expected: true,
		},
		{
			desc:     "Invalid URL - wrong path",
			url:      "https://testartist.bandcamp.com/invalid/path",
			expected: false,
		},
		{
			desc:     "Invalid URL - missing path",
			url:      "https://testartist.bandcamp.com",
			expected: false,
		},
		{
			desc:     "Invalid URL - wrong domain",
			url:      "https://testartist.wrongdomain.com/music",
			expected: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.desc, func() {
			u, err := url.Parse(tc.url)
			s.Require().NoError(err)
			result := isValidBandcampURL(u)
			s.Equal(tc.expected, result)
		})
	}
}

func (s *HandlerTestSuite) Test_matchURLPattern() {
	testCases := []struct {
		desc     string
		url      string
		pattern  string
		expected bool
	}{
		{
			desc:     "Match discography URL",
			url:      "https://testartist.bandcamp.com/music",
			pattern:  "https://{artist}.bandcamp.com/music",
			expected: true,
		},
		{
			desc:     "Match album URL",
			url:      "https://testartist.bandcamp.com/album/testalbum",
			pattern:  "https://{artist}.bandcamp.com/album/{album}",
			expected: true,
		},
		{
			desc:     "Match track URL",
			url:      "https://testartist.bandcamp.com/track/testtrack",
			pattern:  "https://{artist}.bandcamp.com/track/{track}",
			expected: true,
		},
		{
			desc:     "No match - wrong path",
			url:      "https://testartist.bandcamp.com/wrong/path",
			pattern:  "https://{artist}.bandcamp.com/music",
			expected: false,
		},
		{
			desc:     "No match - extra path component",
			url:      "https://testartist.bandcamp.com/album/testalbum/extra",
			pattern:  "https://{artist}.bandcamp.com/album/{album}",
			expected: false,
		},
		{
			desc:     "No match - wrong domain",
			url:      "https://testartist.wrongdomain.com/music",
			pattern:  "https://{artist}.bandcamp.com/music",
			expected: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.desc, func() {
			result := matchURLPattern(tc.url, tc.pattern)
			s.Equal(tc.expected, result)
		})
	}
}

func (s *HandlerTestSuite) Test_Scrapp() {
	testCases := []struct {
		desc           string
		url            string
		scrapper       *scrapper.MockScrapper
		expectedStatus int
		expectedBody   string
	}{
		{
			desc:           "Valid discography URL",
			url:            "https://testartist.bandcamp.com/music",
			scrapper:       s.mockDiscographyScrapper,
			expectedStatus: http.StatusOK,
			expectedBody:   "Request processed successfully",
		},
		{
			desc:           "Valid album URL",
			url:            "https://testartist.bandcamp.com/album/testalbum",
			scrapper:       s.mockAlbumScrapper,
			expectedStatus: http.StatusOK,
			expectedBody:   "Request processed successfully",
		},
		{
			desc:           "Valid track URL",
			url:            "https://testartist.bandcamp.com/track/testtrack",
			scrapper:       s.mockTrackScrapper,
			expectedStatus: http.StatusOK,
			expectedBody:   "Request processed successfully",
		},
		{
			desc:           "Invalid URL - wrong domain",
			url:            "https://testartist.wrongdomain.com/music",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid Bandcamp URL",
		},
		{
			desc:           "Invalid URL - wrong path",
			url:            "https://testartist.bandcamp.com/wrong/path",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid Bandcamp URL",
		},
		{
			desc:           "Empty URL",
			url:            "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Scrapp url param is required",
		},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest("GET", "/api/v1/scrapp", nil)
		s.Require().NoError(err)
		q := req.URL.Query()
		q.Add("url", tc.url)
		req.URL.RawQuery = q.Encode()

		if tc.scrapper != nil {
			tc.scrapper.EXPECT().Execute(gomock.Any()).Return(nil).AnyTimes()
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(s.handler.Scrapp)

		handler.ServeHTTP(rr, req)

		s.Equal(tc.expectedStatus, rr.Code)
		s.Equal(tc.expectedBody, strings.TrimSpace(rr.Body.String()))
	}
}
