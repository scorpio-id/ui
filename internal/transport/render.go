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

	return &WebRender{
		templates: tmpl,
	}
}


func (wr *WebRender) LandingPageHandler(w http.ResponseWriter, r *http.Request) {	
	// website elements that are reused
	components := struct {}{}
		
	// set header
	w.Header().Set("Content-type", "text/html")

	err := wr.templates.ExecuteTemplate(w, "landingPanel.html", components)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}