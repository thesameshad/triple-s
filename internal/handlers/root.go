package handlers

import (
	"net/http"
	"strings"
)

func Root(w http.ResponseWriter, r *http.Request) {
	apiURL := r.URL.Path
	parts := strings.Split(strings.Trim(apiURL, "/"), "/")
	if len(parts) > 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bag Request"))
		return
	}
	switch r.Method {
	case http.MethodGet:

	case http.MethodPut:
	case http.MethodDelete:
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server is running"))
}
