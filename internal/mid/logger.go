package mid

import (
	"github.com/harrisonbrock/gargesale/internal/platform/web"
	"log"
	"net/http"
	"time"
)

func Logger(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		h := func(w http.ResponseWriter, r *http.Request) error {

			start := time.Now()

			// Run the handler chain and catch any propagated error.
			err := before(w, r)

			//Log
			log.Printf(
				"%s %s (%v)",
				r.Method, r.URL.Path,
				time.Since(start),
			)
			// Return nil to indicate the error has been handled.
			return err
		}

		return h
	}

	return f
}
