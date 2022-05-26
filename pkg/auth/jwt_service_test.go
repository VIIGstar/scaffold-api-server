package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var id = "1"
var walletAddress = "0xFC7C98fF48Aa50D75b77A3CA3E7f528817b88255"
var urlRequest, _ = url.Parse(fmt.Sprintf("localhost:8080/v1/auth?wallet_address=%v", walletAddress))

var ginCtx = &gin.Context{
	Request: &defaultRequest,
}
var defaultRequest = http.Request{
	Method: http.MethodPost,
	URL:    urlRequest,
	Header: map[string][]string{
		"User-Agent": {"PostmanRuntime/7.29.0"},
	},
	Body:       nil,
	RemoteAddr: "[::1]:60844",
}

func GetSampleJWTService() *jwtService {
	sv := newJWTService()
	sv.secretKey = "12345678"
	sv.secretKeyBytes = []byte(sv.secretKey)
	return sv
}

func defaultExpired() time.Time {
	timestamp := int64(1650537398)
	return time.Unix(timestamp, 0)
}

func expectToken() string {
	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiIweEZDN0M5OGZGNDhBYTUwRDc1Yjc3QTNDQTNFN2Y1Mjg4MTdiODgyNTUiLCJhdWQiOlsiWzo6MV06NjA4NDQiLCJQb3N0bWFuUnVudGltZS83LjI5LjAiXSwiZXhwIjoxNjUwNTM3Mzk4LCJqdGkiOiIxIn0.PAJYsb33ejQpZPWJMBGYskt9OFXIX06BKDHgeWLh1s8"
}

func TestJwtService_ParseClaims(t *testing.T) {
	claims := parseClaimString(ginCtx)
	assert.Equal(t, claims[0], defaultRequest.RemoteAddr)
	assert.Equal(t, claims[1], defaultRequest.UserAgent())
}

func TestJwtService_GenerateToken(t *testing.T) {
	sv := GetSampleJWTService()
	token := sv.generateToken(walletAddress, id, ginCtx, defaultExpired())
	assert.Equal(t, expectToken(), token)
}

func TestJwtService_ValidateToken(t *testing.T) {
	sv := GetSampleJWTService()
	token := sv.generateToken(walletAddress, id, ginCtx, defaultExpired())
	tk, err := sv.validateToken(token, ginCtx)
	assert.Equal(t, err.Error(), ErrTokenExpired.Error())
	assert.NotNil(t, tk)
}

func TestJwtService_ValidateToken_ExpectOK(t *testing.T) {
	sv := GetSampleJWTService()
	token := sv.generateToken(walletAddress, id, ginCtx, time.Now().Add(time.Hour))

	rq := defaultRequest
	rq.Header.Set("Authorization", token)

	tk, err := sv.validateToken(token, &gin.Context{Request: &rq})
	assert.Nil(t, err)
	assert.True(t, tk.Valid)
	assert.NotNil(t, tk)
	claims := tk.Claims.(jwt.MapClaims)
	assert.Equal(t, claims[IssuerAddressClaimKey], walletAddress)
	assert.Equal(t, claims[IssuerIdClaimKey], id)
}

func TestJwtService_ValidateToken_ExpectFail2(t *testing.T) {
	sv := GetSampleJWTService()
	token := sv.generateToken(walletAddress, id, ginCtx, time.Now().Add(time.Hour))

	rq := defaultRequest
	rq.Header.Set("Authorization", "abc"+token)

	_, err := sv.validateToken("abc"+token, &gin.Context{Request: &rq})
	assert.NotNil(t, err)
	assert.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
}
