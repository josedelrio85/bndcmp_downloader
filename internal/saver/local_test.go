package saver

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestLocalSaver(t *testing.T) {
	suite.Run(t, new(TestLocalSaverSuite))
}

type TestLocalSaverSuite struct {
	suite.Suite
	tempDir string
}

func (s *TestLocalSaverSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "localsaver_test")
	s.Require().NoError(err)
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

func (s *TestLocalSaverSuite) TestLocalSaver_Save() {
	tests := []struct {
		name        string
		folder      *string
		data        string
		filename    string
		expectedErr bool
	}{
		{
			name:        "Successful save",
			folder:      &s.tempDir,
			data:        "Hello, World!",
			filename:    "test.txt",
			expectedErr: false,
		},
		{
			name:        "Save with nil folder",
			folder:      nil,
			data:        "Test content",
			filename:    "test_nil_folder.txt",
			expectedErr: false,
		},
		{
			name:        "Invalid folder",
			folder:      func() *string { s := filepath.Join(s.tempDir, "non_existent"); return &s }(),
			data:        "Test content",
			filename:    "test_invalid_folder.txt",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			saver := NewLocalSaver(tt.folder)
			reader := strings.NewReader(tt.data)

			err := saver.Save(reader, tt.filename)

			if tt.expectedErr {
				s.Error(err)
			} else {
				s.NoError(err)

				// Verify the file was created and contains the correct data
				filePath := tt.filename
				if tt.folder != nil {
					filePath = filepath.Join(*tt.folder, tt.filename)
				}
				content, err := os.ReadFile(filePath)
				s.NoError(err)
				s.Equal(tt.data, string(content))
			}
		})
	}
}

func (s *TestLocalSaverSuite) TestLocalSaver_Save_NonDirectoryFolder() {
	nonDirFile := filepath.Join(s.tempDir, "not_a_directory")
	err := os.WriteFile(nonDirFile, []byte(""), 0644)
	s.Require().NoError(err)

	saver := NewLocalSaver(&nonDirFile)
	err = saver.Save(strings.NewReader("test"), "test.txt")

	s.Error(err)
	s.Equal(os.ErrNotExist, err)
}

func (s *TestLocalSaverSuite) TestLocalSaver_Save_IOError() {
	saver := NewLocalSaver(&s.tempDir)
	errorReader := &errorReader{err: io.ErrUnexpectedEOF}

	err := saver.Save(errorReader, "test_io_error.txt")

	s.Error(err)
	s.Equal(io.ErrUnexpectedEOF, err)
}

type errorReader struct {
	err error
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}
