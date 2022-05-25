package config

import (
	"log"
	"text/template"

	"github.com/Talha2299/bookings/internal/models"
	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the app config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
