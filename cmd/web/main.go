package main

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
	"github.com/sophiabrandt/go-party-finder/business/app/apc"
	"github.com/sophiabrandt/go-party-finder/business/app/handlers"
	"github.com/sophiabrandt/go-party-finder/business/app/session"
	"github.com/sophiabrandt/go-party-finder/config"
	"github.com/sophiabrandt/go-party-finder/foundation/database"
	"github.com/sophiabrandt/go-party-finder/foundation/web"
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
	// Start Database

	log.Println("main: Initializing database support")

	db, err := database.Open(database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}

	defer func() {
		log.Printf("main: Database Stopping : %s", cfg.DB.Host)
		db.Close()
	}()

	// =========================================================================
	// Start Debug Service
	//
	// /debug/pprof - Added to the default mux by importing the net/http/pprof package.
	// /debug/vars - Added to the default mux by importing the expvar package.
	//
	// Not concerned with shutting this down when the application is shutdown.

	log.Println("main: Initializing debugging support")

	go func() {
		log.Printf("main: Debug Listening %s", cfg.Server.DebugAddr)
		if err := http.ListenAndServe(cfg.Server.DebugAddr, http.DefaultServeMux); err != nil {
			log.Printf("main: Debug Listener closed : %v", err)
		}
	}()

	// =========================================================================
	// App Starting

	expvar.NewString("build").Set(build)
	log.Printf("main : Started : Application initializing : version %q", build)
	defer log.Println("main: Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		log.Fatal(err, "generating config for output")
	}
	log.Printf("main: config :\n%v\n", out)

	// =========================================================================
	// TemplateCache, Sessions

	web.NewTemplates(&cfg)
	tc, err := web.CreateTemplateCache()
	if err != nil {
		return errors.Wrap(err, "cannot create template cache")
	}

	session.NewSession(&cfg)
	ses := session.NewStore()

	// put session store and template cache on app context
	apc := apc.New(ses, tc)

	// =========================================================================
	// Start Server

	log.Println("main: Initializing server")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// pass location of static files to router
	staticFilesDir := cfg.Web.StaticFilesLocation

	s := http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      handlers.Router(build, shutdown, apc, staticFilesDir, log, db),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main: APP listening on %s", s.Addr)
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
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and shed load.
		if err := s.Shutdown(ctx); err != nil {
			s.Close()
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
