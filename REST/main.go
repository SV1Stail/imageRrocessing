package main

import (
	"net/http"
	"os"

	"github.com/SV1Stail/imageRrocessing/REST/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	RestPort string
)

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	RestPort = os.Getenv("REST_CONTAINER_PORT")
}

func main() {
	if err := server.InitGRPCClient(); err != nil {
		log.Err(err).Msg("REST server stopped")
		os.Exit(1)
	}

	http.HandleFunc("/upload", server.UploadHandler)
	fs := http.FileServer(http.Dir("/web_ui"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "/web_ui/ui.html")
			return
		}
		fs.ServeHTTP(w, r)
	})

	log.Info().Msgf("HTTP server running on :%s...", RestPort)
	err := http.ListenAndServe(":"+RestPort, nil)
	log.Err(err).Msg("HTTP server fall")
}
