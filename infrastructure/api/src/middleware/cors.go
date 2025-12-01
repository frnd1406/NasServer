package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nas-ai/api/src/config"
	"github.com/sirupsen/logrus"
)

func CORS(cfg *config.Config, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		// Dynamic Origin Handling for LAN/Dev Access
		if origin != "" {
			reqHeaders := c.Request.Header.Get("Access-Control-Request-Headers")
			allowHeaders := "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, x-csrf-token, Authorization, accept, origin, Cache-Control, X-Requested-With"
			if reqHeaders != "" {
				allowHeaders = allowHeaders + ", " + reqHeaders
			}
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", allowHeaders)
			c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
			c.Header("Access-Control-Max-Age", "86400")
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
