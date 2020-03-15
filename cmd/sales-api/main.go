package main

import (
	"context"
	"encoding/json"
	"github.com/harrisonbrock/gargesale/internal/platform/database"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// =========================================================================
	// App Starting

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	// =========================================================================
	// Setup Dependencies
	db, err := database.Open()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// =========================================================================
	// Start API Service

	ps := ProductService{db: db}
	api := http.Server{
		Addr:         "localhost:8000",
		Handler:      http.HandlerFunc(ps.List),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main : API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("error: listening and serving: %s", err)

	case <-shutdown:
		log.Println("main : Start shutdown")

		// Give outstanding requests a deadline for completion.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
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

// Product is something we sale
type Product struct {
	ID          string    `db:"product_id" json:"id"`
	Name        string    `json:"name"`
	Cost        int       `json:"cost"`
	Quantity    int       `json:"quantity"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

// ProductService has handler methods for dealing with Products.
type ProductService struct {
	db *sqlx.DB
}

// ListProduct is a basic HTTP Handler.
func (p *ProductService) List(w http.ResponseWriter, r *http.Request) {

	// Create a slice of products.
	list := []Product{}

	const q = `SELECT * FROM products`
	if err := p.db.Select(&list, q); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error querying data source", err)
		return
	}

	data, err := json.Marshal(list)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error marshalling", err)
		return
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Println("error writing", err)
	}

}
