package handler

import (
	"avito2024/pkg"
	"github.com/gin-gonic/gin"
)

func (h *Handler) SetRequestID(c *gin.Context) {
	requestID := pkg.UUID()
	c.Set("requestID", requestID)
	c.Next()
}
