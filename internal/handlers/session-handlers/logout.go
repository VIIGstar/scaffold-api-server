package session_handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *sessionHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Live!",
	})
}
