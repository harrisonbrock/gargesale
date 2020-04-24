package web

import (
	"context"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"time"
)

// Handler in signature that all applications handler will implement.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entry-point for all web apps.
type App struct {
	mux *chi.Mux
	log *log.Logger
	mw  []Middleware
}

type ctxKey int

const KeyValues ctxKey = 1

type Values struct {
	StatusCode int
	Start      time.Time
}

// NewApp knows how to construct a internal state for App.
func NewApp(logger *log.Logger, mw ...Middleware) *App {
	return &App{
		mux: chi.NewRouter(),
		log: logger,
		mw:  mw,
	}
}

func (a *App) Handle(method, pattern string, h Handler, mw ...Middleware) {

	// First wrap handler specific middleware around this handler.
	h = wrapMiddleware(mw, h)

	// Add the application's general middleware to the handler chain.
	h = wrapMiddleware(a.mw, h)

	fn := func(w http.ResponseWriter, r *http.Request) {

		v := Values{
			StatusCode: 0,
			Start:      time.Now(),
		}
		ctx := context.WithValue(r.Context(), KeyValues, &v)

		if err := h(ctx, w, r); err != nil {

			a.log.Printf("Unhandled error: %+v", err)
		}
	}
	a.mux.MethodFunc(method, pattern, fn)
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
