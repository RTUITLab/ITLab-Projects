package githubreq

import (
	"errors"
)

var (
	UnexpectedCode = errors.New("UnexpectedCode")
	ErrGetLastPage = errors.New("Can't get last page")
)