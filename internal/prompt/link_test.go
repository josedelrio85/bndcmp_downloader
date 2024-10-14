package prompt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/josedelrio85/bndcmp_downloader/internal/scrapper"
	"github.com/stretchr/testify/suite"
)

type TestLinkSuite struct {
	suite.Suite
	controller   *gomock.Controller
	mockPrompter *MockStringPrompter
}

func TestLink(t *testing.T) {
	suite.Run(t, new(TestLinkSuite))
}

func (s *TestLinkSuite) SetupTest() {
	s.controller = gomock.NewController(s.T())
	s.mockPrompter = NewMockStringPrompter(s.controller)
}

func (s *TestLinkSuite) TearDownTest() {
	s.controller.Finish()
}

func (s *TestLinkSuite) TestScrapTypeQuestionLink_Handle() {
	tests := []struct {
		name           string
		input          string
		expectedOutput scrapper.ScrapType
	}{
		{"Track", "1", scrapper.Track},
		{"Album", "2", scrapper.Album},
		{"Discography", "3", scrapper.Discography},
		{"Invalid", "4", scrapper.ScrapType(0)},
	}

	for _, tt := range tests {

		link := NewScrapTypeQuestionLink()
		message := &ChainMessage{}
		link.prompter = s.mockPrompter

		s.mockPrompter.EXPECT().Prompt(gomock.Any()).Return(tt.input)

		link.Handle(message)

		s.Equal(tt.expectedOutput, message.ScrapType)
	}
}

func (s *TestLinkSuite) TestURLCheckerLink_Handle() {
	tests := []struct {
		name        string
		input       string
		scrapType   scrapper.ScrapType
		expectedURL string
		expectError bool
	}{
		{
			name:        "Valid Track URL",
			input:       "https://example.bandcamp.com/track/example-track",
			scrapType:   scrapper.Track,
			expectedURL: "https://example.bandcamp.com/track/example-track",
			expectError: false,
		},
		{
			name:        "Invalid URL",
			input:       "not-a-url",
			scrapType:   scrapper.Track,
			expectedURL: "",
			expectError: true,
		},
		{
			name:        "Mismatched ScrapType",
			input:       "https://example.bandcamp.com/album/example-album",
			scrapType:   scrapper.Track,
			expectedURL: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		link := NewURLCheckerLink()
		message := &ChainMessage{ScrapType: tt.scrapType}
		link.prompter = s.mockPrompter

		s.mockPrompter.EXPECT().Prompt(gomock.Any()).Return(tt.input)

		link.Handle(message)

		if tt.expectError {
			s.Empty(message.URL.Value)
		} else {
			s.Equal(tt.expectedURL, message.URL.Value)
		}
	}
}

func (s *TestLinkSuite) TestStorageQuestionLink_Handle() {
	question := `Where do you want to save the files?
	1. Current directory
	2. Custom directory`

	tests := []struct {
		name               string
		input              string
		mockPrompterOutput string
		expectedOutput     string
	}{
		{"Current Directory", question, "1", "."},
		{"Custom Directory", question, "2", "downloads"},
		{"Invalid Directory", question, "0", ""},
	}

	for _, tt := range tests {
		link := NewStorageQuestionLink()
		message := &ChainMessage{}
		link.prompter = s.mockPrompter

		s.mockPrompter.EXPECT().Prompt(gomock.Any()).Return(tt.mockPrompterOutput)
		if tt.mockPrompterOutput == "2" {
			s.mockPrompter.EXPECT().Prompt(gomock.Any()).Return(tt.expectedOutput)
		}

		link.Handle(message)

		s.Equal(tt.expectedOutput, message.StorageType)
	}
}

func (s *TestLinkSuite) TestSetNext() {
	scrapTypeLink := NewScrapTypeQuestionLink()
	urlCheckerLink := NewURLCheckerLink()
	storageLink := NewStorageQuestionLink()

	scrapTypeLink.SetNext(urlCheckerLink)
	urlCheckerLink.SetNext(storageLink)

	s.Equal(urlCheckerLink, scrapTypeLink.next)
	s.Equal(storageLink, urlCheckerLink.next)
	s.Nil(storageLink.next)
}
