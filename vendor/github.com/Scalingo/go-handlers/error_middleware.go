package handlers

import (
	"encoding/json"
	stderr "errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Scalingo/go-utils/errors"
	"github.com/Scalingo/go-utils/logger"
	"github.com/codegangsta/negroni"
	"github.com/sirupsen/logrus"
)

var ErrorMiddleware MiddlewareFunc = MiddlewareFunc(func(handler HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		log, ok := r.Context().Value("logger").(logrus.FieldLogger)
		if !ok {
			log = logrus.New()
		}

		defer func() {
			if rec := recover(); rec != nil {
				debug.PrintStack()
				err, ok := rec.(error)
				if !ok {
					err = stderr.New(rec.(string))
				}
				log.WithError(err).Error("recover panic")
				w.WriteHeader(500)
				fmt.Fprintln(w, err)
			}
		}()

		rw := negroni.NewResponseWriter(w)
		err := handler(rw, r, vars)

		if ctxerr, ok := err.(errors.ErrCtx); ok {
			log = logger.Get(ctxerr.Ctx())
		}

		if err != nil {
			log.WithField("error", err).Error("request error")
			writeError(rw, err)
		}

		return err
	}
})

func writeError(w negroni.ResponseWriter, err error) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "text/plain")
	}

	// If the status is 0, In means WriteHeader has not been called
	// and we've to write it, otherwise it has been done in the handler
	// with another response code.
	if w.Status() == 0 {
		w.WriteHeader(500)
	}

	if w.Header().Get("Content-Type") == "application/json" {
		json.NewEncoder(w).Encode(&(map[string]string{"error": err.Error()}))
	} else {
		fmt.Fprintln(w, err)
	}
}
