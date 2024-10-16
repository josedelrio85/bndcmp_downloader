package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/josedelrio85/bndcmp_downloader/internal/handler"
	"github.com/josedelrio85/bndcmp_downloader/internal/scrapper"
	"github.com/josedelrio85/bndcmp_downloader/internal/setup"
	"github.com/rs/cors"
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

func setupRouter(httpHandler *handler.HttpHandler) http.Handler {
	r := mux.NewRouter()
	apiV1 := r.PathPrefix("/api/v1").Subrouter()
	apiV1.HandleFunc("/health", httpHandler.Health).Methods("GET")
	// apiV1.HandleFunc("/{artist}", httpHandler.GetDiscography).Methods("GET")
	// apiV1.HandleFunc("/{artist}/{album}", httpHandler.GetAlbum).Methods("GET")
	// apiV1.HandleFunc("/{artist}/track/{track}", httpHandler.GetTrack).Methods("GET")
	apiV1.HandleFunc("/scrapp", httpHandler.Scrapp).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173", "http://localhost:8080", "http://192.168.50.84:8080"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})
	return c.Handler(r)
}
