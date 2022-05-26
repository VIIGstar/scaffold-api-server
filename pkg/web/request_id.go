package web

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	contextRequestID string = "contextRequestID"
	headerRequestID  string = "X-Request-ID"
)

// GetRequestID retrieves request id from context
func GetRequestID(ctx *gin.Context) string {
	if rId, ok := ctx.Get(contextRequestID); ok {
		return rId.(string)
	}
	return ""
}

// RequestIDMiddleware is a middleware that injects UUID to each web request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		attachRequestId(c)
		c.Writer.Header().Set(headerRequestID, GetRequestID(c))
		c.Next()
	}
}

func attachRequestId(ctx *gin.Context) {
	if rId := ctx.Request.Header.Get(headerRequestID); len(rId) > 0 {
		ctx.Set(contextRequestID, rId)
		return
	}
	ctx.Set(contextRequestID, uuid.New().String())
}
