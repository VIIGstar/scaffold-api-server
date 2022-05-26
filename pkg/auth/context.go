package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	pkgerrors "github.com/pkg/errors"
)

var (
	ErrUserNotFoundInCtx = errors.New("user not found in context")
)

const (
	userKey = "user"
)

type UserDetails struct {
	UserID  int64
	Address string
}

// setUserContext sets user details in context
func setUserContext(c *gin.Context, details UserDetails) {
	c.Set(userKey, details)
}

// GetUserFromContext gets user details from context
func GetUserFromContext(c *gin.Context) (UserDetails, error) {
	ud, ok := c.Get(userKey)
	if !ok {
		return UserDetails{}, pkgerrors.WithStack(ErrUserNotFoundInCtx)
	}
	return ud.(UserDetails), nil
}
