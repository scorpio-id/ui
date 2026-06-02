package transport

import (
	"log"
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/scorpio-id/ui/internal/config"
	"github.com/scorpio-id/ui/internal/data"
	"github.com/scorpio-id/ui/internal/tls"
)

func NewRouter(cfg config.Config) *mux.Router {
	// create gorilla mux router
	router := mux.NewRouter()

	persistClient := data.NewPersistenceClient(cfg)

	// create file server for css
	fs := http.FileServer(http.Dir("internal/resources"))

	// strip the prefix
	stripped := http.StripPrefix("/resources/", fs)

	// add stripped file server contents to gorilla mux path prefix handler
	router.PathPrefix("/resources/").Handler(stripped).Methods(http.MethodGet)

	render := NewWebRender(cfg)

	// handles main sections of website
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	}).Methods(http.MethodGet)

	router.HandleFunc("/dashboard", render.DashboardPageHandler).Methods(http.MethodGet)

	// Dashboard routes
	router.HandleFunc("/ui/dashboard/overview", render.ServePartialHandler("overview.html")).Methods(http.MethodGet)

	// OAuth2 routes
	router.HandleFunc("/ui/oauth2", render.ServePartialHandler("oauth2.html")).Methods(http.MethodGet)

	// PKI routes
	router.HandleFunc("/ui/pki", render.ServePartialHandler("pki.html")).Methods(http.MethodGet)

	// Kerberos routes
	router.HandleFunc("/ui/kerberos", render.ServePartialHandler("kerberos.html")).Methods(http.MethodGet)

	// User management routes
	router.HandleFunc("/ui/settings", render.ServePartialHandler("settings.html")).Methods(http.MethodGet)

	// OAuth2 metadata endpoints
	router.HandleFunc("/ui/metadata", render.HandleOAuth2Metadata).Methods(http.MethodGet)
	router.HandleFunc("/ui/register", render.HandleOAuth2Register).Methods(http.MethodPost)

		// check if TLS is enabled, if so create cert client and serialize x509 if on linux OS
	if runtime.GOOS == "linux" {
		// TODO replace with granter.ObtainWebServerIdentity()
		content, err := persistClient.LoadWebPKCS12()
		if err != nil {
			log.Fatal(err)
		}

		// serialize PKCS12 for SSL
		err = tls.SerializePKCS12(content, "/etc/ssl/certs")
		if err != nil {
			log.Fatal(err)
		}
	}

	return router
}
