package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
	"github.com/sophiabrandt/go-party-finder/internal/config"
	"github.com/sophiabrandt/go-party-finder/internal/render"
	"github.com/sophiabrandt/go-party-finder/internal/router"
)

// build is the git version of this program. It is set using build flags in the makefile.
var (
	build = "develop"
	cfg   config.Conf
)

func main() {
	log := log.New(os.Stdout, "PARTYFINDER : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	if err := run(log); err != nil {
		log.Println("main: error:", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {
	// =========================================================================
	// Parse Configuration

	cfg.Version.SVN = build
	cfg.Version.Desc = "Apache 2.0 License"

	if err := conf.Parse(os.Args[1:], "PARTYFINDER", &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage("PARTYFINDER", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		case conf.ErrVersionWanted:
			version, err := conf.VersionString("PARTYFINDER", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// =========================================================================
	// App Starting

	log.Printf("main : Started : Application initializing : version %q", build)
	defer log.Println("main: Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		log.Fatal(err, "generating config for output")
	}
	log.Printf("main: config :\n%v\n", out)

	// =========================================================================
	// Create TemplateCache

	tc, err := render.CreateTemplateCache()
	if err != nil {
		return errors.Wrap(err, "cannot create template cache")
	}
	cfg.App.TemplateCache = tc
	render.NewTemplates(&cfg)

	// =========================================================================
	// Start Server

	log.Println("main: Initializing server")

	s := http.Server{
		Addr:         cfg.Web.Addr,
		Handler:      router.New(log),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server startup failed")
	}

	return nil
}
