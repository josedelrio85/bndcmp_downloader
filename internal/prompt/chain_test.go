package prompt

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestChain(t *testing.T) {
	suite.Run(t, new(ChainTestSuite))
}

type ChainTestSuite struct {
	suite.Suite
	controller *gomock.Controller
	links      []Link
}

func (c *ChainTestSuite) SetupTest() {
	c.controller = gomock.NewController(c.T())
}

func (c *ChainTestSuite) TearDownTest() {
	c.controller.Finish()
}

func (c *ChainTestSuite) TestNewChain() {
	link_one := NewMockLink(c.controller)
	link_two := NewMockLink(c.controller)
	link_three := NewMockLink(c.controller)

	c.links = []Link{link_one, link_two, link_three}

	// Expect SetNext to be called with the next link
	link_one.EXPECT().SetNext(link_two)
	link_two.EXPECT().SetNext(link_three)

	// Expect Handle to be called on the first link
	link_one.EXPECT().Handle(gomock.Any())

	// Create the chain
	chain := NewChain(c.links)

	c.Assert().NotNil(chain)
	c.Assert().Equal(3, len(chain.Links))
	c.Assert().Equal(c.links, chain.Links)
}
