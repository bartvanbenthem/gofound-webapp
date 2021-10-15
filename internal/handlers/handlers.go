package handlers

import (
	"log"
	"net/http"

	"github.com/bartvanbenthem/gofound-webapp/internal/config"
	"github.com/bartvanbenthem/gofound-webapp/internal/forms"
	"github.com/bartvanbenthem/gofound-webapp/internal/models"
	"github.com/bartvanbenthem/gofound-webapp/internal/render"
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

// TestForm is the handler for the testform page
func (m *Repository) TestForm(w http.ResponseWriter, r *http.Request) {
	var emptyTestForm models.TestForm
	data := make(map[string]interface{})
	data["testform"] = emptyTestForm

	render.RenderTemplate(w, r, "testform.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostTestForm handles post
func (m *Repository) PostTestForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	tf := models.TestForm{
		Name:  r.Form.Get("name"),
		Email: r.Form.Get("email"),
	}

	form := forms.New(r.PostForm)

	form.Required("name", "email")
	form.MinLength("name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["testform"] = tf
		render.RenderTemplate(w, r, "testform.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
}
