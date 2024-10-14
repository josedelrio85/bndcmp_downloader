package prompt

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/josedelrio85/bndcmp_downloader/internal/bandcamp"
	"github.com/josedelrio85/bndcmp_downloader/internal/scrapper"
)

//go:generate mockgen -source=$GOFILE -package $GOPACKAGE -destination link_mock.go
type Link interface {
	Handle(message *ChainMessage)
	SetNext(next Link) Link
}

type ScrapTypeQuestionLink struct {
	next     Link
	prompter StringPrompter
}

func NewScrapTypeQuestionLink() *ScrapTypeQuestionLink {
	return &ScrapTypeQuestionLink{
		prompter: &DefaultStringPrompter{},
	}
}

func (r *ScrapTypeQuestionLink) Handle(message *ChainMessage) {
	question := `What do you want to download?
	1. Track
	2. Album
	3. Discography
	`
	inputValue := r.prompter.Prompt(question)

	scrapTypes := map[string]scrapper.ScrapType{
		"1": scrapper.Track,
		"2": scrapper.Album,
		"3": scrapper.Discography,
	}

	if scrapType, ok := scrapTypes[inputValue]; ok {
		message.ScrapType = scrapType
		if r.next != nil {
			r.next.Handle(message)
		}
	} else {
		log.Println("Invalid value")
	}
}

func (r *ScrapTypeQuestionLink) SetNext(next Link) Link {
	r.next = next
	return r
}

type URLCheckerLink struct {
	next     Link
	prompter StringPrompter
}

func NewURLCheckerLink() *URLCheckerLink {
	return &URLCheckerLink{
		prompter: &DefaultStringPrompter{},
	}
}

func (c *URLCheckerLink) Handle(message *ChainMessage) {
	url := c.prompter.Prompt(c.getQuestion(message))
	bandcampURL, err := c.processBandcampURL(url, message.ScrapType)
	if err != nil {
		log.Printf("Error processing URL: %v\n", err)
		return
	}

	message.URL = bandcampURL

	if c.next != nil {
		c.next.Handle(message)
	}
}

func (c *URLCheckerLink) getQuestion(message *ChainMessage) string {
	comment := ""
	if message.ScrapType == scrapper.Track {
		comment = strings.Join([]string{comment, "\n", "Example: https://{band-name}.bandcamp.com/track/{track-name}"}, " ")
	} else if message.ScrapType == scrapper.Album {
		comment = strings.Join([]string{comment, "\n", "Example: https://{band-name}.bandcamp.com/album/{album-name}"}, " ")
	} else if message.ScrapType == scrapper.Discography {
		comment = strings.Join([]string{comment, "\n", "Example: https://{band-name}.bandcamp.com/music"}, " ")
	}
	return strings.Join([]string{"Enter the URL: ", comment}, " ")
}

func (c *URLCheckerLink) processBandcampURL(url string, expectedScrapType scrapper.ScrapType) (bandcamp.BandcampURL, error) {
	bandcampURL := bandcamp.BandcampURL{Value: url}

	if err := bandcampURL.Parse(); err != nil {
		return bandcamp.BandcampURL{}, fmt.Errorf("parsing URL: %w", err)
	}

	if err := bandcampURL.Validate(); err != nil {
		return bandcamp.BandcampURL{}, fmt.Errorf("validating URL: %w", err)
	}

	scrapTypeValue := bandcampURL.Classify()
	if expectedScrapType != scrapper.ScrapType(scrapTypeValue) {
		return bandcamp.BandcampURL{}, fmt.Errorf("invalid URL: expected scrap type %v, got %v", expectedScrapType, scrapTypeValue)
	}

	return bandcampURL, nil
}

func (c *URLCheckerLink) SetNext(next Link) Link {
	c.next = next
	return c
}

type StorageQuestionLink struct {
	next     Link
	prompter StringPrompter
}

func NewStorageQuestionLink() *StorageQuestionLink {
	return &StorageQuestionLink{
		prompter: &DefaultStringPrompter{},
	}
}

func (s *StorageQuestionLink) Handle(message *ChainMessage) {
	options := map[string]string{
		"1": "Current directory",
		"2": "Custom directory",
	}

	question := "Where do you want to save the files?\n"
	for key, value := range options {
		question += fmt.Sprintf("\t%s. %s\n", key, value)
	}

	storageType := s.prompter.Prompt(question)

	switch storageType {
	case "1":
		message.StorageType = "."
	case "2":
		message.StorageType = s.prompter.Prompt("Enter the custom directory: ")
	default:
		log.Println("Invalid value")
		return
	}

	if s.next != nil {
		s.next.Handle(message)
	}
}

func (s *StorageQuestionLink) SetNext(next Link) Link {
	s.next = next
	return next
}

type StringPrompter interface {
	Prompt(label string) string
}

type DefaultStringPrompter struct{}

// Prompt asks for a string value using the label
func (d *DefaultStringPrompter) Prompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}
