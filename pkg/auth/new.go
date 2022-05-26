package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cast"
	"logur.dev/logur"
	"net/http"
	info_log "scaffold-api-server/pkg/info-log"
	"strings"
	"time"
)

type Authenticator interface {
	GenerateAccessToken(address, id string, c *gin.Context) string
}

type impl struct {
}

func (i impl) GenerateAccessToken(address, id string, c *gin.Context) string {
	return newJWTService().generateToken(address, id, c, time.Now().Add(AccessTokenExpiry))
}

func New() Authenticator {
	return impl{}
}

// Middleware checks if request comes with valid access token
func Middleware(logger logur.LoggerFacade) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !skipTokenCheck(c.Request.URL.Path) {
			at := extractAccessToken(c)
			if at == "" {
				logger.Info("[Auth] no access token, proceed to extract refresh token")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// if access token present
			t, err := newJWTService().validateToken(at, c)
			if err != nil {
				logger.Error("[Auth] error when verify access token",
					info_log.ErrorToLogFields("details: ", err))
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// Extract token metadata and store its tokenDetails into the same request context
			userDetails, err := extractTokenMetadata(t)
			if err != nil || userDetails == nil {
				logger.Error("[Auth] error when extract token metadata",
					info_log.ErrorToLogFields("details: ", err))
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			setUserContext(c, *userDetails)
			c.Next()
		}
	}
}

func extractAccessToken(c *gin.Context) string {
	atCookie, err := c.Request.Cookie(AccessTokenKey)
	if err != nil || len(atCookie.Value) == 0 {
		at := c.GetHeader(strings.ToUpper(AccessTokenKey))
		if len(at) == 0 {
			return ""
		}

		return at
	}
	return atCookie.Value
}

// extractTokenMetadata extracts metadata from token
func extractTokenMetadata(token *jwt.Token) (*UserDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrMapClaimNotFound
	}

	userID, ok := claims[IssuerIdClaimKey]
	if !ok {
		return nil, ErrUserIDNotFound
	}

	address, ok := claims[IssuerAddressClaimKey]
	if !ok {
		return nil, ErrUserAddressNotFound
	}

	return &UserDetails{
		UserID:  cast.ToInt64(userID),
		Address: address.(string),
	}, nil
}

func skipTokenCheck(uri string) bool {
	skipPath := []string{
		"/api/v1/liveness",
		"/api/v1/readiness",
		"/api/v1/sessions/login",
		"/api/v1/investors/signup",
		"/api/v1/debug/pprof",
		"/swagger",
	}
	for _, s := range skipPath {
		if strings.Contains(uri, s) {
			return true
		}
	}
	return false
}
