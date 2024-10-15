package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/josedelrio85/bndcmp_downloader/internal/scrapper"
)

type HttpHandler struct {
	baseFolder          string
	discographyScrapper scrapper.Scrapper
	albumScrapper       scrapper.Scrapper
	trackScrapper       scrapper.Scrapper
}

func NewHttpHandler(
	baseFolder string,
	discographyScrapper scrapper.Scrapper,
	albumScrapper scrapper.Scrapper,
	trackScrapper scrapper.Scrapper,
) *HttpHandler {
	return &HttpHandler{
		baseFolder:          baseFolder,
		discographyScrapper: discographyScrapper,
		albumScrapper:       albumScrapper,
		trackScrapper:       trackScrapper,
	}
}

func (h *HttpHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *HttpHandler) GetDiscography(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetDiscography")
	vars := mux.Vars(r)
	artist := vars["artist"]
	if artist == "" {
		http.Error(w, "Artist is required", http.StatusBadRequest)
		return
	}
	discographyLink := fmt.Sprintf("https://%s.bandcamp.com/music", artist)
	discographyURL, err := url.Parse(discographyLink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.discographyScrapper.Execute(discographyURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request processed of " + artist + " downloaded successfully"))
}

func (h *HttpHandler) GetAlbum(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetAlbum")
	vars := mux.Vars(r)
	artist := vars["artist"]
	if artist == "" {
		http.Error(w, "Artist is required", http.StatusBadRequest)
		return
	}
	album := vars["album"]
	if album == "" {
		http.Error(w, "Album is required", http.StatusBadRequest)
		return
	}
	albumLink := fmt.Sprintf("https://%s.bandcamp.com/album/%s", artist, album)
	albumURL, err := url.Parse(albumLink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.albumScrapper.Execute(albumURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Album " + album + " of " + artist + " downloaded successfully"))
}

func (h *HttpHandler) GetTrack(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetTrack")
	vars := mux.Vars(r)
	artist := vars["artist"]
	if artist == "" {
		http.Error(w, "Artist is required", http.StatusBadRequest)
		return
	}
	track := vars["track"]
	if track == "" {
		http.Error(w, "Track is required", http.StatusBadRequest)
		return
	}
	trackLink := fmt.Sprintf("https://%s.bandcamp.com/track/%s", artist, track)
	trackURL, err := url.Parse(trackLink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.trackScrapper.Execute(trackURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Track " + track + " of " + artist + " downloaded successfully"))
}

func (h *HttpHandler) Scrapp(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	scrapParam := queryParams.Get("url")
	log.Println("Scrapping url: ", scrapParam)
	if scrapParam == "" {
		http.Error(w, "Scrapp url param is required", http.StatusBadRequest)
		return
	}

	scrapURL, err := url.Parse(scrapParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !isValidBandcampURL(scrapURL) {
		http.Error(w, "Invalid Bandcamp URL", http.StatusBadRequest)
		return
	}

	scrapper, err := h.getScrapper(scrapURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = scrapper.Execute(scrapURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request processed successfully"))
}

func (h *HttpHandler) getScrapper(scrapURL *url.URL) (scrapper.Scrapper, error) {
	path := scrapURL.Path
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		return nil, fmt.Errorf("invalid URL")
	}

	if pathParts[1] == "music" {
		return h.discographyScrapper, nil
	} else if pathParts[1] == "album" {
		return h.albumScrapper, nil
	} else if pathParts[1] == "track" {
		return h.trackScrapper, nil
	}
	return nil, fmt.Errorf("invalid scrap type")
}

func isValidBandcampURL(u *url.URL) bool {
	validPatterns := []string{
		"https://{artist}.bandcamp.com/music",
		"https://{artist}.bandcamp.com/album/{album}",
		"https://{artist}.bandcamp.com/track/{track}",
	}

	for _, pattern := range validPatterns {
		if matchURLPattern(u.String(), pattern) {
			return true
		}
	}
	return false
}

func matchURLPattern(url, pattern string) bool {
	regexPattern := strings.ReplaceAll(pattern, ".", "\\.")
	regexPattern = strings.ReplaceAll(regexPattern, "{artist}", "[^/]+")
	regexPattern = strings.ReplaceAll(regexPattern, "{album}", "[^/]+")
	regexPattern = strings.ReplaceAll(regexPattern, "{track}", "[^/]+")
	regexPattern = "^" + regexPattern + "$"

	match, _ := regexp.MatchString(regexPattern, url)
	return match
}
