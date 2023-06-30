package service

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func writeJSON(w http.ResponseWriter, v any) {
	content, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if _, err := w.Write(content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func readJSON(r *http.Request, v any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := r.Body.Close(); err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}
