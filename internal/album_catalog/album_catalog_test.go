package album_catalog

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AlbumCatalogTestSuite struct {
	suite.Suite
	tempDir string
	catalog *InMemoryAlbumCatalog
}

func TestAlbumCatalogSuite(t *testing.T) {
	suite.Run(t, new(AlbumCatalogTestSuite))
}

func (s *AlbumCatalogTestSuite) SetupTest() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "album_catalog_test")
	s.Require().NoError(err)
	s.catalog = NewInMemoryAlbumCatalog(s.tempDir)
}

func (s *AlbumCatalogTestSuite) TearDownTest() {
	os.RemoveAll(s.tempDir)
}

func (s *AlbumCatalogTestSuite) TestGenerate_EmptyDirectory() {
	err := s.catalog.Generate()
	s.Require().NoError(err)
	s.Empty(s.catalog.MapDir)
}

func (s *AlbumCatalogTestSuite) TestGenerate_SingleFile() {
	filename := "test.txt"
	err := os.WriteFile(filepath.Join(s.tempDir, filename), []byte("test"), 0644)
	s.Require().NoError(err)

	err = s.catalog.Generate()
	s.Require().NoError(err)
	s.Len(s.catalog.MapDir, 1)
	s.True(s.catalog.MapDir[filename])
}

func (s *AlbumCatalogTestSuite) TestGenerate_NestedDirectories() {
	nestedDir := filepath.Join(s.tempDir, "nested")
	err := os.Mkdir(nestedDir, 0755)
	s.Require().NoError(err)

	filename1 := "test1.txt"
	filename2 := "test2.txt"
	err = os.WriteFile(filepath.Join(s.tempDir, filename1), []byte("test1"), 0644)
	s.Require().NoError(err)
	err = os.WriteFile(filepath.Join(nestedDir, filename2), []byte("test2"), 0644)
	s.Require().NoError(err)

	err = s.catalog.Generate()
	s.Require().NoError(err)
	s.Len(s.catalog.MapDir, 3) // 2 files + 1 directory
	s.True(s.catalog.MapDir[filename1])
	s.True(s.catalog.MapDir[filename2])
	s.True(s.catalog.MapDir["nested"])
}

func (s *AlbumCatalogTestSuite) TestGenerate_NonExistentDirectory() {
	s.catalog.baseFolder = "/non/existent/directory"
	err := s.catalog.Generate()
	s.Require().Error(err)
}

func (s *AlbumCatalogTestSuite) TestGenerate_InvalidDirectory() {
	s.catalog.baseFolder = "/dev/null"
	err := s.catalog.Generate()
	s.Require().Error(err)
}
