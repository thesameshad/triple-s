package server

import (
	"fmt"
	"net/http"
	"triple-s/internal/handlers"
)

func Routes() *http.ServeMux {
	fmt.Println("-------------------Server Started--------------------")
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Root)
	mux.HandleFunc("/health", handlers.Health)
	return mux
}
