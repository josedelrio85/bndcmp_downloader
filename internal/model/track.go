package model

type Track struct {
	Title       string
	TrackNumber int64
	Artist      string
	Album       *string
	URL         string
	DownloadURL string
}
