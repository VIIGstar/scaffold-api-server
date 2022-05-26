package session_handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	query_params "scaffold-api-server/internal/query-params"
	"scaffold-api-server/pkg/auth"
	info_log "scaffold-api-server/pkg/info-log"
)

// @Summary  Validate user then get access token
// @Tags     Session
// @Param    wallet_address  query  string  true  "public key address to user wallet"
// @Accept   json
// @Produce  json
// @Success  200      {object}  auth.Authentication
// @Failure  400,500  {object}  object{error=string}
// @Router   /api/v1/sessions/login [post]
func (h *sessionHandler) Login(c *gin.Context) {
	investor, err := h.repo.Database().User().Find(c, query_params.GetUserParams{
		Address: c.Query("wallet_address"),
	}, false)
	if err != nil {
		h.logger.Error("error find user", info_log.ErrorToLogFields("details", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, auth.Authentication{
		AccessToken: auth.New().GenerateAccessToken(investor.Address, fmt.Sprintf("%v", investor.ID), c),
		Success:     true,
	})
}
