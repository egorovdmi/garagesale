package web

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func Respond(w http.ResponseWriter, value interface{}, statusCode int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "marshaling error")
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "error writing")
	}

	return nil
}

func RespondError(w http.ResponseWriter, err error) error {
	if webErr, ok := err.(*Error); ok {
		er := ErrorResponse{
			Error: webErr.Error(),
		}

		return Respond(w, er, webErr.Status)
	}

	resp := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}

	return Respond(w, resp, http.StatusInternalServerError)
}
