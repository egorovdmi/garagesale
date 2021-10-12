package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/egorovdmi/garagesale/internal/platform/web"
	"github.com/egorovdmi/garagesale/internal/product"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Product struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// ProductsList returns the list with all products
func (ps *Product) List(w http.ResponseWriter, r *http.Request) error {

	list, err := product.List(r.Context(), ps.DB)
	if err != nil {
		return errors.Wrap(err, "db select error")
	}

	if err := web.Respond(w, list, http.StatusOK); err != nil {
		return errors.Wrap(err, "error writing")
	}

	return nil
}

// Single returns the product by id
func (ps *Product) Single(w http.ResponseWriter, r *http.Request) error {

	id := web.GetURLParam(r, "id")
	p, err := product.Single(r.Context(), ps.DB, id)
	if err != nil {
		switch err {
		case product.ErrInvalidId:
			return web.NewRequestError(err, http.StatusBadRequest)
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrap(err, "db select error")
		}
	}

	if err := web.Respond(w, p, http.StatusOK); err != nil {
		return errors.Wrap(err, "error writing")
	}

	return nil
}

// Single returns the product by id
func (ps *Product) Create(w http.ResponseWriter, r *http.Request) error {

	var np product.NewProduct
	if err := web.Decode(r, &np); err != nil {
		return err
	}

	p, err := product.Create(r.Context(), ps.DB, np, time.Now())
	if err != nil {
		return errors.Wrap(err, "db create error")
	}

	if err := web.Respond(w, p, http.StatusCreated); err != nil {
		return errors.Wrap(err, "error writing")
	}

	return nil
}
