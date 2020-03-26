package mid

import (
	"github.com/harrisonbrock/gargesale/internal/platform/web"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"time"
)

func Logger(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		h := func(w http.ResponseWriter, r *http.Request) error {

			v, ok := r.Context().Value(web.KeyValues).(*web.Values)

			if !ok {
				return errors.New("web values missing from context")
			}

			// Run the handler chain and catch any propagated error.
			err := before(w, r)

			//Log
			log.Printf(
				"%d %s %s (%v)",
				v.StatusCode,
				r.Method, r.URL.Path,
				time.Since(v.Start),
			)
			// Return nil to indicate the error has been handled.
			return err
		}

		return h
	}

	return f
}
