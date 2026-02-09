package server

import (
	"fmt"
	"net/http"
	"triple-s/internal/handlers"
)

var started bool

func StartServer(port int) error {
	addr := fmt.Sprintf(":%d", port)
	if !started {
		handlers.SetDataDir(dataDir)
		started = true
	}
	return http.ListenAndServe(addr, Routes())
}

var dataDir string

func SetDir(dir string) { dataDir = dir }
