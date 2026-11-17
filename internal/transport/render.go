package transport

import (

	"log"
	"net/http"
	"text/template"

	"github.com/scorpio-id/ui/internal/config"
)

type WebRender struct {
	templates *template.Template
}

func NewWebRender(cfg config.Config) *WebRender {
	// TODO: add template directory to config
	// Matches and stores all templates found in the templates section
	// template files are stored as their file name within the resource folder
	tmpl, err := template.ParseGlob("./internal/resources/templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// Also parse partial templates
	tmpl, err = tmpl.ParseGlob("./internal/resources/templates/partials/*.html")
	if err != nil {
		log.Fatal(err)
	}

	return &WebRender{
		templates: tmpl,
	}
}

func (wr *WebRender) DashboardPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	err := wr.templates.ExecuteTemplate(w, "dashboardPanel.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ServePartialHandler serves individual panel partials
func (wr *WebRender) ServePartialHandler(templateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html; charset=utf-8")
		err := wr.templates.ExecuteTemplate(w, templateName, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
func (wr *WebRender) HandleOAuth2Metadata(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{
		"issuer": "https://oauth.scorpio.ordinarycomputing.com:8082",
		"authorization_endpoint": "https://oauth.scorpio.ordinarycomputing.com:8082/authorize",
		"token_endpoint": "https://oauth.scorpio.ordinarycomputing.com:8082/token",
		"userinfo_endpoint": "https://oauth.scorpio.ordinarycomputing.com:8082/userinfo",
		"jwks_uri": "https://oauth.scorpio.ordinarycomputing.com:8082/.well-known/jwks.json"
	}`))
}

// HandleLogout handles user logout
func (wr *WebRender) HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Clear session or token here
	// For now, just redirect to login or home
	http.Redirect(w, r, "/", http.StatusFound)
}


// HandleOAuth2Register handles OAuth2 client registration
func (wr *WebRender) HandleOAuth2Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	clientID := r.FormValue("client_id")
	email := r.FormValue("email")

	// Process registration here
	_ = clientID
	_ = email

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status": "success", "message": "Application registered successfully"}`))
}
