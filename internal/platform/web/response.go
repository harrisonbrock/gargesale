package web

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

func Response(w http.ResponseWriter, val interface{}, status int) error {

	data, err := json.Marshal(val)

	if err != nil {
		return errors.Wrap(err, "marshaling value to json")
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "writing to client")
	}

	return nil
}

func RespondError(w http.ResponseWriter, err error) error {

	if webErr, ok := err.(*Error); ok {
		resp := ErrorResponse{Error: webErr.Err.Error()}
		return Response(w, resp, webErr.Status)
	}

	resp := ErrorResponse{Error: http.StatusText(http.StatusInternalServerError)}
	return Response(w, resp, http.StatusInternalServerError)
}
