package main

import (
	"log"
	"runtime"
	"net/http"

	"github.com/scorpio-id/ui/internal/config"
	"github.com/scorpio-id/ui/internal/transport"
)

func main() {
	// parse local config
	cfg := config.NewConfig("internal/config/local.yml")

	// create a new mux router
	router := transport.NewRouter(cfg)

	// start the server
	if runtime.GOOS == "linux" {
		log.Fatal(http.ListenAndServeTLS(":"+cfg.Server.Port, "/etc/ssl/certs/scorpio-ui.pem", "/etc/ssl/certs/scorpio-ui.key", router))
	} else {
		log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, router))
	}
}