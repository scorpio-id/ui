package transport

import (
	"crypto/rsa"
	"encoding/json"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

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

// HandleOAuth2Clients fetches OAuth2 clients metadata and returns as HTML table rows
func (wr *WebRender) HandleOAuth2Clients(w http.ResponseWriter, r *http.Request) {
	// Fetch clients data from external OAuth2 server
	resp, err := http.Get("https://oauth.scorpio.ordinarycomputing.com:8082/ui/metadata")
	if err != nil {
		http.Error(w, "Failed to fetch clients data", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusBadGateway)
		return
	}

	// Parse the JSON response
	var data struct {
		Clients []struct {
			ID                   string          `json:"id"`
			Email                string          `json:"email"`
			UserPrincipalName    string          `json:"user_principal_name"`
			ServicePrincipalName string          `json:"service_principal_name"`
			CommonName           string          `json:"common_name"`
			SubjectAltNames      []string        `json:"subject_alternate_names"`
			Authorizations       json.RawMessage `json:"authorizations"`
		} `json:"clients"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Failed to parse clients data", http.StatusBadGateway)
		return
	}

	// Render HTML table rows
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	for _, client := range data.Clients {
		altNames := strings.Join(client.SubjectAltNames, ", ")
		html := `<tr>
			<td>` + client.ID + `</td>
			<td>` + client.Email + `</td>
			<td>` + client.UserPrincipalName + `</td>
			<td>` + client.ServicePrincipalName + `</td>
			<td>` + client.CommonName + `</td>
			<td>` + altNames + `</td>
		</tr>`
		w.Write([]byte(html))
	}

	// If no clients, show empty state
	if len(data.Clients) == 0 {
		w.Write([]byte(`<tr><td colspan="6">No clients found</td></tr>`))
	}
}

// HandleCertificates fetches certificate metadata and returns as HTML table rows
func (wr *WebRender) HandleCertificates(w http.ResponseWriter, r *http.Request) {
	// Fetch certificate data from external CA server
	resp, err := http.Get("https://ca.scorpio.ordinarycomputing.com:8081/ui/metadata")
	if err != nil {
		http.Error(w, "Failed to fetch certificate data", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusBadGateway)
		return
	}

	// Parse the JSON response
	var data struct {
		Certificates []struct {
			CommonName             string         `json:"common_name"`
			SubjectAlternateNames  []string       `json:"subject_alternate_names"`
			SerialNumber           *big.Int       `json:"serial_number"`
			PublicKey              *rsa.PublicKey `json:"public_key"`
			IssuedDate             time.Time      `json:"issued"`
			ExpirationDate         time.Time      `json:"expires"`
			IsCertificateAuthority bool           `json:"is_certificate_authority"`
		} `json:"certificates"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Failed to parse certificate data", http.StatusBadGateway)
		return
	}

	// Render HTML table rows
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	for _, cert := range data.Certificates {
		altNames := strings.Join(cert.SubjectAlternateNames, ", ")
		html := `<tr>
			<td>` + cert.CommonName + `</td>
			<td>` + altNames + `</td>
			<td>` + cert.SerialNumber.String() + `</td>
			<td>` + cert.ExpirationDate.Format("2006-01-02") + `</td>
			<td>` + strconv.FormatBool(cert.IsCertificateAuthority) + `</td>
		</tr>`
		w.Write([]byte(html))
	}

	// If no certificates, show empty state
	if len(data.Certificates) == 0 {
		w.Write([]byte(`<tr><td colspan="6">No certificates found</td></tr>`))
	}
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
