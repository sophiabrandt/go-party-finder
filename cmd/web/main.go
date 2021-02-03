package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
	"github.com/sophiabrandt/go-party-finder/internal/config"
	"github.com/sophiabrandt/go-party-finder/internal/handlers"
	"github.com/sophiabrandt/go-party-finder/internal/web"
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

	tc, err := web.CreateTemplateCache()
	if err != nil {
		return errors.Wrap(err, "cannot create template cache")
	}
	cfg.App.TemplateCache = tc
	web.NewTemplates(&cfg)

	// =========================================================================
	// Start Server

	log.Println("main: Initializing server")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	s := http.Server{
		Addr:         cfg.Web.Addr,
		Handler:      handlers.Router(build, shutdown, log),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main: API listening on %s", s.Addr)
		serverErrors <- s.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case sig := <-shutdown:
		log.Printf("main: %v : Start shutdown", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := s.Shutdown(ctx); err != nil {
			s.Close()
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
