package libcamera

import (
	"errors"
)

var (
	ErrNoBytes = errors.New("libcamera error: no bytes received from channel!")
)
