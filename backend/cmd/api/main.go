package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/josedelrio85/bndcmp_downloader/internal/handler"
	"github.com/josedelrio85/bndcmp_downloader/internal/scrapper"
	"github.com/josedelrio85/bndcmp_downloader/internal/setup"
)

func main() {
	log.Println("Starting Bandcamp downloader API")

	httpHandler := setupHttpHHandler()
	router := setupRouter(httpHandler)

	// Start the HTTP server
	addr := ":8099"
	log.Printf("Server is listening on port %s...\n", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Println(err)
	}
}

func setupHttpHHandler() *handler.HttpHandler {
	config := setup.LoadConfig()

	discographyScrapper := scrapper.NewDiscographyScrapper(config.Retriever, config.Parser, config.Saver, config.AlbumCatalog)
	albumScrapper := scrapper.NewAlbumScrapper(config.Retriever, config.Parser, config.Saver, config.AlbumCatalog)
	trackScrapper := scrapper.NewTrackScrapper(config.Retriever, config.Parser, config.Saver, config.AlbumCatalog)

	return handler.NewHttpHandler(
		config.BaseFolder,
		discographyScrapper,
		albumScrapper,
		trackScrapper,
	)
}

func setupRouter(httpHandler *handler.HttpHandler) *mux.Router {
	r := mux.NewRouter()
	apiV1 := r.PathPrefix("/api/v1").Subrouter()
	apiV1.HandleFunc("/health", httpHandler.Health).Methods("GET")
	apiV1.HandleFunc("/{artist}", httpHandler.GetDiscography).Methods("GET")
	apiV1.HandleFunc("/{artist}/{album}", httpHandler.GetAlbum).Methods("GET")
	apiV1.HandleFunc("/{artist}/track/{track}", httpHandler.GetTrack).Methods("GET")
	return r
}