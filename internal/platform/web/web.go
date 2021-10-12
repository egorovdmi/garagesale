package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// App is an entry point for all web applications
type App struct {
	mux *chi.Mux
	log *log.Logger
}

type Handler func(http.ResponseWriter, *http.Request) error

// NewApp knows how to construct internal state of an App
func NewApp(log *log.Logger) *App {
	return &App{
		mux: chi.NewRouter(),
		log: log,
	}
}

func (app *App) Handle(method string, pattern string, handler Handler) {
	app.mux.MethodFunc(method, pattern, func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			app.log.Printf("ERROR : %v\n", err)
			if err := RespondError(w, err); err != nil {
				app.log.Printf("ERROR : %v\n", err)
			}
		}
	})
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.mux.ServeHTTP(w, r)
}
