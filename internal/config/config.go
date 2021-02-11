package config

import (
	"html/template"
	"time"

	"github.com/golangcollege/sessions"
	"github.com/ardanlabs/conf"
)

// Conf holds all the app configuration.
type Conf struct {
	conf.Version
	Web struct {
		Addr            string        `conf:"default:0.0.0.0:8000"`
		DebugAddr       string        `conf:"default:0.0.0.0:6060"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:5s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
		IdleTimeout     time.Duration `conf:"default:120s"`
	}
	DB struct {
		User       string `conf:"default:postgres"`
		Password   string `conf:"default:postgres,noprint"`
		Host       string `conf:"default:0.0.0.0:8461"`
		Name       string `conf:"default:postgres"`
		DisableTLS bool   `conf:"default:true"`
	}
	App struct {
		UseCache            bool   `conf:"default:false"`
		StaticFilesLocation string `conf:"default:./ui/static"`
		TemplateLocation    string `conf:"default:./ui/html"`
		TemplateCache       map[string]*template.Template
		InProduction        bool   `conf:"default:false"`
		SessionSecret       string `conf:"default:UzXRAAKMfSBWvV44HN25tRhfpKFMKT7a"`
	}
}

// LocalContext holds local context, e.g., template cache, session manager.
type LocalContext struct {
	TemplateCache map[string]*template.Template
	Session       *sessions.Session
}
