package server

import (
	"fmt"
	"net/http"
)

func StartServer(port int) error {
	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, Routes())
}
