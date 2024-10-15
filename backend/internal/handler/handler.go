package handler

import (
	"fmt"
	"net/http"
	"net/url"

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
	w.Write([]byte("Discography of " + artist + " downloaded successfully"))
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
