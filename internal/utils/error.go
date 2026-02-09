package utils

import "errors"

var (
	ErrInvalidPort = errors.New("port must be in range 1024–65535")
	ErrInvalidDir  = errors.New("directory must be an absolute path")
	ErrCreateDir   = errors.New("directory creation failed")
)
