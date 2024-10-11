package scrapper

type ScrapType int

const (
	Undefined ScrapType = iota
	Track
	Album
	Discography
)
