package saver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/josedelrio85/bndcmp_downloader/internal/model"
	"github.com/stretchr/testify/suite"
)

func TestLocalSaver(t *testing.T) {
	suite.Run(t, new(TestLocalSaverSuite))
}

type TestLocalSaverSuite struct {
	suite.Suite
	saver   *LocalSaver
	tempDir string
}

func (s *TestLocalSaverSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "localsaver_test")
	s.Require().NoError(err)

	s.saver = NewLocalSaver(&s.tempDir)
}

func (s *TestLocalSaverSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
	// Check for and remove test_nil_folder.txt in the current directory
	nilFolderTestFile := "test_nil_folder.txt"
	if _, err := os.Stat(nilFolderTestFile); err == nil {
		err := os.Remove(nilFolderTestFile)
		s.Require().NoError(err, "Failed to remove test_nil_folder.txt")
	}
}

func toPointer(s string) *string {
	return &s
}

func (s *TestLocalSaverSuite) Test_generateDirectoryStructure() {
	testCases := []struct {
		Description string
		Track       *model.Track
		Expected    string
	}{
		{
			Description: "No album",
			Track: &model.Track{
				Title:       "Elbow",
				TrackNumber: 1,
				Artist:      "King Gizzard and the Lizard Wizard",
			},
			Expected: "King Gizzard and the Lizard Wizard",
		},
		{
			Description: "With album",
			Track: &model.Track{
				Title:       "Elbow",
				Artist:      "King Gizzard and the Lizard Wizard",
				TrackNumber: 1,
				Album:       toPointer("12 Bar Bruise"),
			},
			Expected: "King Gizzard and the Lizard Wizard/12 Bar Bruise",
		},
		{
			Description: "No track number",
			Track: &model.Track{
				Title:  "Elbow",
				Artist: "King Gizzard and the Lizard Wizard",
			},
			Expected: "King Gizzard and the Lizard Wizard",
		},
	}

	for _, tt := range testCases {
		filename := s.saver.generateDirectoryStructure(tt.Track)
		s.Equal(tt.Expected, filename)
	}
}

func (s *TestLocalSaverSuite) Test_checkFolder() {
	testCases := []struct {
		name        string
		folderPath  string
		expectedErr bool
		setup       func(string) error
		teardown    func(string) error
	}{
		{
			name:        "Existing folder",
			folderPath:  s.tempDir,
			expectedErr: false,
		},
		{
			name:        "Non-existing folder",
			folderPath:  filepath.Join(s.tempDir, "new_folder"),
			expectedErr: false,
			teardown: func(path string) error {
				return os.RemoveAll(path)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.setup != nil {
				err := tc.setup(tc.folderPath)
				s.Require().NoError(err)
			}

			err := s.saver.checkFolder(tc.folderPath)

			if tc.expectedErr {
				s.Error(err)
			} else {
				s.NoError(err)
				_, err := os.Stat(tc.folderPath)
				s.NoError(err, "Folder should exist")
			}

			if tc.teardown != nil {
				err := tc.teardown(tc.folderPath)
				s.Require().NoError(err)
			}
		})
	}
}

func (s *TestLocalSaverSuite) Test_saveFile() {
	testCases := []struct {
		Description    string
		Filename       string
		ExpectedResult bool
	}{
		{
			Description:    "Success",
			Filename:       "test.txt",
			ExpectedResult: true,
		},
		{
			Description:    "Fail",
			Filename:       "foobar/test.txt",
			ExpectedResult: false,
		},
	}

	for _, tt := range testCases {
		fmt.Println(tt.Description)
		err := s.saver.saveFile(s.tempDir, tt.Filename, strings.NewReader("test"))
		if tt.ExpectedResult {
			s.NoError(err)
		} else {
			s.Error(err)
		}
	}
}

func (s *TestLocalSaverSuite) TestSave() {
	tests := []struct {
		name        string
		track       *model.Track
		data        string
		expectedErr bool
	}{
		{
			name: "Save track without album",
			track: &model.Track{
				Title:       "Test Track",
				TrackNumber: 1,
				Artist:      "Test Artist",
			},
			data:        "Test audio data",
			expectedErr: false,
		},
		{
			name: "Save track with album",
			track: &model.Track{
				Title:       "Test Track",
				TrackNumber: 2,
				Artist:      "Test Artist",
				Album:       toPointer("Test Album"),
			},
			data:        "Test audio data with album",
			expectedErr: false,
		},
		{
			name:        "Save with nil track",
			track:       nil,
			data:        "Test data",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			reader := strings.NewReader(tt.data)
			err := s.saver.Save(reader, tt.track)

			if tt.expectedErr {
				s.Error(err)
			} else {
				s.NoError(err)

				// Verify the file was created and contains the correct data
				if tt.track != nil {
					filePath := s.saver.generateDirectoryStructure(tt.track)
					filePath = filepath.Join(s.tempDir, filePath)
					fileName := fmt.Sprintf("%02d - %s.mp3", tt.track.TrackNumber, tt.track.Title)
					filePath = filepath.Join(filePath, fileName)

					// Check if the file exists
					fileInfo, err := os.Stat(filePath)
					s.NoError(err, "File should exist")
					s.False(fileInfo.IsDir(), "File path should not be a directory")

					content, err := os.ReadFile(filePath)
					s.NoError(err, "Should be able to read the file")
					s.Equal(tt.data, string(content), "File content should match the input data")
				}
			}
		})
	}
}
