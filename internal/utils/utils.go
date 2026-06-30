package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Payload map[string]interface {
}

func WriteJSON(w http.ResponseWriter, statusCode int, p Payload) error {
	jsn, err := json.MarshalIndent(p, "", "")

	if err != nil {
		return err
	}

	jsn = append(jsn, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsn)

	return nil
}

func RetrieveIDFromRequest(r *http.Request) (int64, error) {
	idParams := chi.URLParam(r, "id")
	if idParams == "" {
		return 0,
			errors.New("id parameter is missing")
	}

	id, err := strconv.ParseInt(idParams, 10, 64)
	if err != nil {
		return 0,
			errors.New("invalid id parameter")
	}
	return id, nil
}
