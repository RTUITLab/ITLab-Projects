package statuscode

import (
	"errors"
	"fmt"
)

type StatusCode struct {
	Err		error
	Status	int
}

func (s *StatusCode) Error() string {
	return fmt.Sprintf("%v: %v", s.Status, s.Err)
}

func WrapStatusError(
	err 	error,
	status 	int,
) error {
	return &StatusCode{
		Err: err,
		Status: status,
	}
}

func GetStatus(
	err error,
) (status int, ok bool) {
	StatusCode, ok := err.(*StatusCode)
	if !ok {
		return 0, ok
	}

	return StatusCode.Status, ok
}

func GetError(
	err error,
) error {
	StatusCode, ok := err.(*StatusCode)
	if !ok {
		return nil
	}

	return StatusCode.Err
}

func Is(
	err		error,
	target 	error,
) bool {
	if StatusCode, ok := err.(*StatusCode); ok {
		return errors.Is(StatusCode.Err, target)
	} else {
		return false
	}
}