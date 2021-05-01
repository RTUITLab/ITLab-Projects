package mfsreq

import "errors"

var (
	ErrUnexpectedCode 	= errors.New("Unexpected code")
	NetError			= errors.New("Error when send req")
)

type UnexpectedCodeErr struct {
	Err error
	Code int
}

func (uce *UnexpectedCodeErr) Error() string {
	return uce.Err.Error()
}