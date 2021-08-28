package handlers

import (
	"net/http"

	"github.com/bartvanbenthem/gofound-web/pkg/config"
	"github.com/bartvanbenthem/gofound-web/pkg/models"
	"github.com/bartvanbenthem/gofound-web/pkg/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	// send data to the template
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Items is the handler for the items page
func (m *Repository) Items(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "items.page.tmpl", &models.TemplateData{})
}

// Contact is the handler for the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// PostTest handles post
func (m *Repository) PostTest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("succesfull post\n"))
}
