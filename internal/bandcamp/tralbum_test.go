package bandcamp

import (
	"testing"

	"github.com/josedelrio85/bndcmp_downloader/internal/model"
	"github.com/stretchr/testify/suite"
)

func TestTrAlbum(t *testing.T) {
	suite.Run(t, new(TestTrAlbumSuite))
}

type TestTrAlbumSuite struct {
	suite.Suite
}

func (s *TestTrAlbumSuite) TestToTrack() {
	tests := []struct {
		name     string
		trAlbum  *TrAlbum
		expected *model.Track
	}{
		{
			name: "Valid TrAlbum",
			trAlbum: &TrAlbum{
				Current:  Current{Title: "Test Track"},
				Artist:   "Test Artist",
				URL:      "https://example.com/track",
				AlbumURL: "/album/test-album",
				Trackinfo: []TrackInfo{
					{File: File{Mp3128: "https://example.com/download"}},
				},
			},
			expected: &model.Track{
				Title:       "Test Track",
				Artist:      "Test Artist",
				Album:       toPointer("Test Album"),
				URL:         "https://example.com/track",
				DownloadURL: "https://example.com/download",
			},
		},
		{
			name:     "Nil TrAlbum",
			trAlbum:  nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := tt.trAlbum.ToTrack()
			s.Equal(tt.expected, result)
		})
	}
}

func (s *TestTrAlbumSuite) TestGetAlbumName() {
	tests := []struct {
		name     string
		trAlbum  *TrAlbum
		expected *string
	}{
		{
			name: "Valid AlbumURL",
			trAlbum: &TrAlbum{
				AlbumURL: "/album/test-album-name",
			},
			expected: toPointer("Test Album Name"),
		},
		{
			name: "Empty AlbumURL",
			trAlbum: &TrAlbum{
				AlbumURL: "",
			},
			expected: toPointer(""),
		},
		{
			name:     "Nil TrAlbum",
			trAlbum:  nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := tt.trAlbum.getAlbumName()
			s.Equal(tt.expected, result)
		})
	}
}

func toPointer(s string) *string {
	return &s
}
