// This program performs administrative tasks for the garage sale service.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/egorovdmi/garagesale/internal/platform/conf"
	"github.com/egorovdmi/garagesale/internal/platform/database"
	"github.com/egorovdmi/garagesale/internal/schema"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Println("shutting down", "error:", err)
		os.Exit(1)
	}
}

func run() error {

	// =========================================================================
	// Configuration

	var cfg struct {
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:postgres"`
			DisableSSL bool   `conf:"default:false"`
		}
		Args conf.Args
	}

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// Initialize dependencies.
	db, err := database.Open(database.Config{
		User:       cfg.DB.User,
		Pass:       cfg.DB.Password,
		Host:       cfg.DB.Host,
		DBName:     cfg.DB.Name,
		DisableSSL: cfg.DB.DisableSSL,
	})
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer db.Close()

	switch cfg.Args.Num(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			return errors.Wrap(err, "error applying migrations")
		}
		fmt.Println("Migrations complete")
		return nil

	case "seed":
		if err := schema.Seed(db); err != nil {
			return errors.Wrap(err, "error seeding database")
		}
		fmt.Println("Seed data complete")
		return nil
	}

	return nil
}
