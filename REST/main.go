package main

import (
	"net/http"

	"github.com/SV1Stail/imageRrocessing/REST/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

}

func main() {
	server.InitGRPCClient()

	http.HandleFunc("/upload", server.UploadHandler)
	fs := http.FileServer(http.Dir("./web_ui"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "./web_ui/ui.html")
			return
		}
		fs.ServeHTTP(w, r)
	})

	log.Info().Msg("HTTP server running on :8081...")
	err := http.ListenAndServe(":8081", nil)
	log.Err(err).Msg("HTTP server fall")
}
