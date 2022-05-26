package database

import "strings"

func IsDuplicateErr(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry")
}
