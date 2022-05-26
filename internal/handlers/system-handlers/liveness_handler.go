package system_handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Check API still alive
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} object{message=string}
// @Failure 400,500 {object} object{error=string}
// @Router /api/v1/liveness [get]
func (h *systemHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
