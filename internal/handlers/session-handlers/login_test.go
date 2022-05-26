package session_handlers

import (
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"scaffold-api-server/internal/entities"
	"scaffold-api-server/internal/repository"
	"scaffold-api-server/internal/services/database"
	"scaffold-api-server/internal/services/log"
	"scaffold-api-server/pkg/auth"
	"scaffold-api-server/pkg/config"
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
		RemoteAddr: "[::1]:60844",
	}
)

func TestSessionHandler_Login(t *testing.T) {
	sqlDB, mock, _ := sqlmock.New()
	investor := entities.Investor{
		DefaultModel: entities.DefaultModel{
			ID: 1,
		},
		Address:     walletAddress,
		NetworkName: "",
		NetworkURL:  "",
		ChainID:     0,
		Symbol:      "",
		RefreshKey:  "",
	}
	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})
	gormDB, err := gorm.Open(dialector)
	if err != nil {
		panic(err)
	}
	db = database.TestDB(gormDB)

	mock.MatchExpectationsInOrder(false)
	mock.ExpectQuery(
		regexp.QuoteMeta(
			"SELECT * FROM `investors` WHERE address = ? AND `investors`.`deleted_at` IS NULL ORDER BY `investors`.`id` LIMIT 1",
		),
	).
		WithArgs(walletAddress).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "address"}).
				AddRow(investor.ID, walletAddress))
	handler := New(logger, repository.New(logger, db, nil))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &defaultRequest
	handler.Login(c)

	obj := auth.Authentication{}
	json.Unmarshal(w.Body.Bytes(), &obj)
	assert.Equal(t, w.Code, http.StatusOK)
	assert.True(t, obj.Success || obj.Error == "Already registered!")
}
