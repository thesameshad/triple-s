package cmd

import (
	"os"
	"triple-s/internal/utils"
)

func ArgsRoute(port int, dir string, help bool) error {
	if help {
		Help()
		return nil
	}

	if port < 1024 || port > 65535 {
		return utils.ErrInvalidPort
	}

	if dir == "" {
		return utils.ErrInvalidDir
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return utils.ErrCreateDir
	}

	return nil
}
