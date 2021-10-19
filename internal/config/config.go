package config

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
	"github.com/bartvanbenthem/gofound-webapp/internal/models"
)

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
