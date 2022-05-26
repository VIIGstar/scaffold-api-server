package auth

import "errors"

var (
	HTTPUnauthorizedError = errors.New("invalid credential")

	ErrUserIDNotFound      = errors.New("user id not found in token claim")
	ErrUserAddressNotFound = errors.New("user address not found in token claim")
	ErrMapClaimNotFound    = errors.New("map claim not found in token")
	ErrTokenExpired        = errors.New("Token is expired")
	ErrTokenAudience       = errors.New("Token invalid audience")
)
