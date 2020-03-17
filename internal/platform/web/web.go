package web

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

// App is the entry-point for all web apps.
type App struct {
	mux *chi.Mux
	log *log.Logger
}

// NewApp knows how to construct a internal state for App.
func NewApp(logger *log.Logger) *App {
	return &App{
		mux: chi.NewRouter(),
		log: logger,
	}
}

func (a *App) Handle(method, pattern string, fn http.HandlerFunc) {
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
