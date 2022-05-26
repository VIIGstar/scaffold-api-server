package auth

import "time"

const AccessTokenExpiry = time.Hour

const (
	AccessTokenKey        = "access-token"
	IssuerAddressClaimKey = "iss"
	IssuerIdClaimKey      = "jti"
)
