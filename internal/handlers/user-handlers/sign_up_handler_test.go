package user_handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"scaffold-api-server/internal/dtos"
	"scaffold-api-server/internal/entities"
	"scaffold-api-server/internal/repository"
	"scaffold-api-server/internal/services/cache"
	"scaffold-api-server/internal/services/database"
	"scaffold-api-server/internal/services/log"
	"scaffold-api-server/pkg"
	"scaffold-api-server/pkg/auth"
	"scaffold-api-server/pkg/config"
	app_http "scaffold-api-server/pkg/http"
	"testing"
)

const (
	walletAddress = "0xFC7C98fF48Aa50D75b77A3CA3E7f528817b88255"
)

var (
	db             *database.DB
	logger         = log.NewLogger(config.LogConfig{})
	urlRequest, _  = url.Parse(fmt.Sprintf("localhost:8080/v1/auth?wallet_address=%v", walletAddress))
	defaultRequest = http.Request{
		Method: http.MethodPost,
		URL:    urlRequest,
		Header: map[string][]string{
			"User-Agent": {"PostmanRuntime/7.29.0"},
		},
		Body:       io.NopCloser(bytes.NewReader(defaultInvestor())),
		RemoteAddr: "[::1]:60844",
	}
)

func InitServer() userHandler {
	v, f := viper.New(), pflag.NewFlagSet(string(pkg.APIAppName), pflag.ExitOnError)
	cfg := config.New(v, f)
	db = database.New(logger, cfg.Database)

	logger.Info("Initializing redis...")
	c, err := cache.New(context.TODO(),
		fmt.Sprintf("%v:%v",
			viper.GetString("redis.host"),
			viper.GetString("redis.port"),
		),
		logger)
	if err != nil {
		panic(err)
	}
	return New(logger, repository.New(logger, db, c))
}

func TestUserHandler_Signup(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &defaultRequest
	v, f := viper.New(), pflag.NewFlagSet(string(pkg.APIAppName), pflag.ExitOnError)
	config.New(v, f)
	handler := InitServer()
	defer db.Close()

	handler.Signup(c)

	obj := auth.Authentication{}
	json.Unmarshal(w.Body.Bytes(), &obj)
	assert.True(t, obj.Success || obj.Error == "Already registered!")
}

func TestUserHandler_SignupInvalidWalletAddress(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// make wrong address
	investor := dtos.InvestorDTO{
		Investor: entities.Investor{
			Address:     "empty",
			NetworkName: "Avalanche",
			NetworkURL:  viper.GetString("network.url"),
			ChainID:     1379,
			Symbol:      "AVAX",
		},
	}

	d, _ := json.Marshal(investor)

	c.Request = &defaultRequest
	c.Request.Body = io.NopCloser(bytes.NewReader(d))

	v, f := viper.New(), pflag.NewFlagSet(string(pkg.APIAppName), pflag.ExitOnError)
	config.New(v, f)
	handler := InitServer()
	defer db.Close()

	handler.Signup(c)

	obj := auth.Authentication{}
	json.Unmarshal(w.Body.Bytes(), &obj)
	assert.Equal(t, app_http.HTTPBadRequestError, obj.Error)
}

func TestUserHandler_SignupFailEmptyBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	v, f := viper.New(), pflag.NewFlagSet(string(pkg.APIAppName), pflag.ExitOnError)
	config.New(v, f)
	handler := InitServer()
	defer db.Close()

	handler.Signup(c)

	obj := auth.Authentication{}
	json.Unmarshal(w.Body.Bytes(), &obj)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func defaultInvestor() []byte {
	obj := dtos.InvestorDTO{
		Investor: entities.Investor{
			Address:     walletAddress,
			NetworkName: "Avalanche",
			NetworkURL:  viper.GetString("network.url"),
			ChainID:     1379,
			Symbol:      "AVAX",
		},
	}

	d, _ := json.Marshal(obj)
	return d
}
