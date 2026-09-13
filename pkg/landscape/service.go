package landscape

import (
	"errors"

	"resty.dev/v3"
)

type Service struct {
	client *resty.Client
}

var ErrNotImplemented = errors.New("not implemented")
