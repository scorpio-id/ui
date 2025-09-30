package transport

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/scorpio-id/ui/internal/config"
)

func NewRouter(cfg config.Config) *mux.Router {
	// create gorilla mux router
	router := mux.NewRouter()

	// create file server for css
	fs := http.FileServer(http.Dir("internal/resources"))

	// strip the prefix
	stripped := http.StripPrefix("/resources/", fs)

	// add stripped file server contents to gorilla mux path prefix handler
	router.PathPrefix("/resources/").Handler(stripped).Methods(http.MethodGet)

	render := NewWebRender(cfg)

	// handles main sections of website
	router.HandleFunc("/", render.LandingPageHandler).Methods(http.MethodGet, http.MethodOptions)
	
	return router
}