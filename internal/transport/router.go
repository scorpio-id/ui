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
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	}).Methods(http.MethodGet)
	router.HandleFunc("/dashboard", render.DashboardPageHandler).Methods(http.MethodGet)

	// Dashboard routes
	router.HandleFunc("/ui/dashboard/overview", render.ServePartialHandler("overview.html")).Methods(http.MethodGet)

	// OAuth2 routes
	router.HandleFunc("/ui/oauth2", render.ServePartialHandler("clients.html")).Methods(http.MethodGet)

	// PKI routes
	router.HandleFunc("/ui/pki/x509s", render.ServePartialHandler("x509s.html")).Methods(http.MethodGet)

	// Kerberos routes
	router.HandleFunc("/ui/kerberos/tgts", render.ServePartialHandler("tgts.html")).Methods(http.MethodGet)

	// User management routes
	router.HandleFunc("/ui/notifications", render.ServePartialHandler("notifications.html")).Methods(http.MethodGet)
	router.HandleFunc("/ui/settings", render.ServePartialHandler("settings.html")).Methods(http.MethodGet)
	router.HandleFunc("/ui/user/profile", render.ServePartialHandler("profile.html")).Methods(http.MethodGet)
	router.HandleFunc("/ui/user/logout", render.HandleLogout).Methods(http.MethodPost)

	// OAuth2 metadata endpoints
	router.HandleFunc("/ui/metadata", render.HandleOAuth2Metadata).Methods(http.MethodGet)
	router.HandleFunc("/ui/register", render.HandleOAuth2Register).Methods(http.MethodPost)

	return router
}
