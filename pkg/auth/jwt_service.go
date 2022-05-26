package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var secretKey string

func init() {
	if len(secretKey) == 0 {
		secretKey = fmt.Sprintf("%v", time.Now().Unix())
	}
}

type jwtService struct {
	secretKey      string
	secretKeyBytes []byte
}

func newJWTService() *jwtService {
	return &jwtService{
		secretKey:      secretKey,
		secretKeyBytes: []byte(secretKey),
	}
}

func (j *jwtService) generateToken(address, id string, c *gin.Context, expired time.Time) string {
	aud := parseClaimString(c)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:   address,
		Audience: aud,
		ExpiresAt: &jwt.NumericDate{
			Time: expired,
		},
		ID: id,
	})
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (j *jwtService) validateToken(authorization string, c *gin.Context) (token *jwt.Token, err error) {
	token, err = jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
		return j.secretKeyBytes, nil
	})
	if err != nil {
		return
	}

	if !token.Valid {
		err = ErrTokenExpired
		return
	}

	aud := parseClaimString(c)
	claims := token.Claims.(jwt.MapClaims)
	for _, check := range aud {
		if !claims.VerifyAudience(check, true) {
			err = ErrTokenAudience
			return
		}
	}
	return token, nil
}

// parseClaimString return array includes in order by host, user agent
func parseClaimString(c *gin.Context) jwt.ClaimStrings {
	var (
		host      = c.Request.RemoteAddr
		userAgent = c.Request.UserAgent()
	)

	return jwt.ClaimStrings{
		host,
		userAgent,
	}
}
