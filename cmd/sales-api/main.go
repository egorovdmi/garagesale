package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/egorovdmi/garagesale/cmd/sales-api/internal/handlers"
	"github.com/egorovdmi/garagesale/internal/platform/conf"
	"github.com/egorovdmi/garagesale/internal/platform/database"
)

func main() {
	// =========================================================================
	// App starting

	fmt.Println("Starting a web server.")
	defer fmt.Println("Server stopped.")

	var cfg struct {
		Web struct {
			Address         string        `conf:"default:localhost:8000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		DB struct {
			User       string `conf:"default:postgres"`
			Pass       string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:postgres"`
			DisableSSL bool   `conf:"default:false"`
		}
	}

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				log.Fatalf("error : generating config usage : %v", err)
			}
			fmt.Println(usage)
			return
		}
		log.Fatalf("error: parsing config: %s", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		log.Fatalf("error : generating config for output : %v", err)
	}
	log.Printf("main : Config :\n%v\n", out)

	// =========================================================================
	// Setup dependencies
	db, err := database.Open(database.Config{Host: cfg.DB.Host, DBName: cfg.DB.Name, User: cfg.DB.User, Pass: cfg.DB.Pass, DisableSSL: cfg.DB.DisableSSL})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ps := handlers.Product{
		DB: db,
	}

	// =========================================================================
	// Start API service

	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      http.HandlerFunc(ps.List),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Println("listening on " + cfg.Web.Address)
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatal(err)
	case <-shutdown:
		log.Println("Shutting down...")
		timeout := cfg.Web.ShutdownTimeout
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = api.Close()
		}

		if err != nil {
			log.Fatalf("main : could not stop server gracefully : %v", err)
		}
	}
}
