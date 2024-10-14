package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/josedelrio85/bndcmp_downloader/internal/album_catalog"
	"github.com/josedelrio85/bndcmp_downloader/internal/handler"
	"github.com/josedelrio85/bndcmp_downloader/internal/parser"
	"github.com/josedelrio85/bndcmp_downloader/internal/retriever"
	"github.com/josedelrio85/bndcmp_downloader/internal/saver"
	"github.com/josedelrio85/bndcmp_downloader/internal/scrapper"
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
	baseFolder := "downloads"

	retrieverClient := retriever.NewHttpClient()
	parseClient := parser.NewParseClient()
	saveClient := saver.NewLocalSaver(&baseFolder)

	albumCatalog := album_catalog.NewInMemoryAlbumCatalog(baseFolder)
	if err := albumCatalog.Generate(baseFolder); err != nil {
		log.Println("Error generating album catalog: ", err)
		os.Exit(1)
	}

	discographyScrapper := scrapper.NewDiscographyScrapper(retrieverClient, parseClient, saveClient, albumCatalog)
	albumScrapper := scrapper.NewAlbumScrapper(retrieverClient, parseClient, saveClient, albumCatalog)
	trackScrapper := scrapper.NewTrackScrapper(retrieverClient, parseClient, saveClient, albumCatalog)

	return handler.NewHttpHandler(
		baseFolder,
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
