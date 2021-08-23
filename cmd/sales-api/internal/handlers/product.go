package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/egorovdmi/garagesale/internal/product"
	"github.com/jmoiron/sqlx"
)

type Product struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// ProductsList returns the list with all products
func (ps *Product) List(w http.ResponseWriter, r *http.Request) {

	list, err := product.List(ps.DB)
	if err != nil {
		ps.Log.Println("db select error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		ps.Log.Println("marshaling error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		ps.Log.Println("error writing", err)
	}
}

// Single returns the product by id
func (ps *Product) Single(w http.ResponseWriter, r *http.Request) {

	id := "TODO"
	p, err := product.Single(ps.DB, id)
	if err != nil {
		ps.Log.Println("db select error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(p)
	if err != nil {
		ps.Log.Println("marshaling error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		ps.Log.Println("error writing", err)
	}
}
