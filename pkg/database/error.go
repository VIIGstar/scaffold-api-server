package database

import (
	"errors"
)

var (
	InvalidRequestError = errors.New("invalid request")
	NotFoundError       = errors.New("not found")
)

func IsNotFoundRecord(err error) bool {
	return err.Error() == NotFoundError.Error()
}
