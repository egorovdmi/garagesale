package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/egorovdmi/garagesale/internal/product"
	"github.com/jmoiron/sqlx"
)

type Product struct {
	DB *sqlx.DB
}

// ProductsList returns the list with all products
func (ps *Product) List(w http.ResponseWriter, r *http.Request) {

	list, err := product.List(ps.DB)
	if err != nil {
		log.Println("db select error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
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
