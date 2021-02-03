package config

import (
	"html/template"
	"time"

	"github.com/ardanlabs/conf"
)

// Conf holds all the app configuration.
type Conf struct {
	conf.Version
	Web struct {
		Addr            string        `conf:"default:0.0.0.0:8000"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:5s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
		IdleTimeout     time.Duration `conf:"default:120s"`
	}
	App struct {
		UseCache      bool `conf:"default:true"`
		TemplateCache map[string]*template.Template
		InProduction  bool `conf:"default:false"`
	}
}
