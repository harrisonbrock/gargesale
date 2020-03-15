package main

import (
	"fmt"
	"github.com/harrisonbrock/gargesale/internal/platform/conf"
	"github.com/harrisonbrock/gargesale/internal/platform/database"
	"github.com/harrisonbrock/gargesale/internal/schema"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
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
			DisableTLS bool   `conf:"default:false"`
		}
		Args conf.Args
	}

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				return errors.Wrap(err, "main : generating usage")
			}
			fmt.Println(usage)
			return nil
		}
		log.Fatalf("error: parsing config: %s", err)
	}

	// Initialize dependencies.
	db, err := database.Open(database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		return errors.Wrap(err, "opening database")
	}
	defer db.Close()

	switch cfg.Args.Num(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			return errors.Wrap(err, "applying migrations")
		}
		fmt.Println("Migrations complete")
		return nil

	case "seed":
		if err := schema.Seed(db); err != nil {
			return errors.Wrap(err, "seeding database")
		}
		fmt.Println("Seed data complete")
		return nil
	}
	return nil
}
