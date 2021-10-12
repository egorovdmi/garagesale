package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

func GetURLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func Decode(r *http.Request, dest interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return NewRequestError(errors.Wrap(err, "decoding request body"), http.StatusBadRequest)
	}

	return nil
}
