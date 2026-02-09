package server

import (
	"fmt"
	"net/http"
	"triple-s/internal/handlers"
)

func Routes() *http.ServeMux {
	fmt.Println("-------------------Server Started--------------------")
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("public"))
	mux.Handle("/ui/", http.StripPrefix("/ui/", fs))
	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/", handlers.Root)
	return mux
}
