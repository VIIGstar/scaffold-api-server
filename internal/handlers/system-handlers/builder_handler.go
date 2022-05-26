package system_handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Check Image build info
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} object{revision=string}
// @Failure 400,500 {object} object{error=string}
// @Router /api/v1/builder [get]
func (h *systemHandler) GetBuildInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":  "ok",
		"revision": h.buildInfo.CommitHash,
	})
}
