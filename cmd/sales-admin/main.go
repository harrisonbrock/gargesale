package main

import (
	"flag"
	"github.com/harrisonbrock/gargesale/internal/platform/database"
	"github.com/harrisonbrock/gargesale/internal/schema"
	_ "github.com/lib/pq"
	"log"
)

func main() {

	// Setup Dependencies
	db, err := database.Open()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	flag.Parse()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Fatal("Applying migrations", err)
		}
		log.Println("Migrations complete")
		return
	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Fatal("Applying seed data", err)
		}
		log.Println("Seed data inserted")
		return
	}
}
