package web

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

// Handler in signature that all applications handler will implement.
type Handler func(w http.ResponseWriter, r *http.Request) error

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

func (a *App) Handle(method, pattern string, h Handler) {

	fn := func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			resp := ErrorResponse{Error: err.Error()}
			if err := Response(w, resp, http.StatusInternalServerError); err != nil {
				a.log.Println(err)
			}
		}
	}
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
