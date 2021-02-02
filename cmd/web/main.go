package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ardanlabs/conf"
)

func main() {
	// =========================================================================
	// Configuration

	var cfg struct {
		Web struct {
			Addr            string        `conf:"default:0.0.0.0:8000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
		}
	}

	if err := conf.Parse(os.Args[1:], "PARTYFINDER", &cfg); err != nil {
		switch err {
		case conf.ErrHelpWanted:
			usage, err := conf.Usage("PARTYFINDER", &cfg)
			if err != nil {
				log.Fatal(err, "generating config usage")
			}
			fmt.Println(usage)
		case conf.ErrVersionWanted:
			version, err := conf.VersionString("PARTYFINDER", &cfg)
			if err != nil {
				log.Fatal(err, "generating config version")
			}
			fmt.Println(version)
		}
		log.Fatal(err, "parsing config")
	}

	// =========================================================================
	// App Starting

	log.Println("main : Started : Application initializing")
	defer log.Println("main: Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		log.Fatal(err, "generating config for output")
	}
	log.Printf("main: Config :\n%v\n", out)

	// =========================================================================
	// Start Server

	log.Println("main: Initializing server")

	mux := http.NewServeMux()
	mux.HandleFunc("/", Greet)

	s := http.Server{
		Addr:         cfg.Web.Addr,
		Handler:      mux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server startup failed")
	}
}

func Greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
