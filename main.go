package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// =========================================================================
	// App starting

	fmt.Println("Starting a web server.")
	defer fmt.Println("Server stopped.")

	// =========================================================================
	// Setup dependencies
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// =========================================================================
	// Start API service

	api := http.Server{
		Addr:         "localhost:8000",
		Handler:      http.HandlerFunc(ProductsList),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Println("listening on localhost:8000")
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatal(err)
	case <-shutdown:
		log.Println("Shutting down...")
		timeout := 5 * time.Second
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

func openDB() (*sqlx.DB, error) {
	q := url.Values{}
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "postgres"),
		Host:     "localhost",
		Path:     "postgres",
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}

type Product struct {
	Name     string `json:"name"`
	Cost     int    `json:"cost"`
	Quantity int    `json:"quantity"`
}

// ProductsList returns the list with all products
func ProductsList(w http.ResponseWriter, r *http.Request) {
	list := []Product{}

	if true {
		list = append(list, Product{Name: "Comic Book", Cost: 99, Quantity: 10})
		list = append(list, Product{Name: "Cyberpunk 2077", Cost: 49, Quantity: 99})
	}

	data, err := json.Marshal(list)
	if err != nil {
		log.Println("marshaling error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Println("error writing", err)
	}
}
