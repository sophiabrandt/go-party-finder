package config

import (
	"time"

	"github.com/ardanlabs/conf"
)

// Conf holds all the app configuration.
type Conf struct {
	conf.Version
	Server struct {
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
	Web struct {
		UseCache            bool   `conf:"default:false"`
		StaticFilesLocation string `conf:"default:./ui/static"`
		TemplateLocation    string `conf:"default:./ui/html"`
		InProduction        bool   `conf:"default:false"`
	}
	Session struct {
		LifeTime time.Duration `conf:"default:12h"`
		Persist  bool          `conf:"default:true"`
		Secret   string        `conf:"u46IpCV9y5Vlur8YvODJEhgOY8m9JVE4"`
	}
}
