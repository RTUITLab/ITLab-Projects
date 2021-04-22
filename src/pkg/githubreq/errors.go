package githubreq

import (
	"errors"
)

var (
	UnexpectedCode = errors.New("UnexpectedCode")
	ErrGetLastPage = errors.New("Can't get last page")
	ErrForbiden = errors.New("Forbidden status from github")
	ErrUnatorizared = errors.New("Unathorizeted status from github")
)