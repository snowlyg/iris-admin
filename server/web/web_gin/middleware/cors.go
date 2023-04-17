package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/snowlyg/iris-admin/server/web"
)

// Cors
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", web.CONFIG.Cors.AccessOrigin)
		c.Header("Access-Control-Allow-Headers", web.CONFIG.Cors.AccessHeaders)
		c.Header("Access-Control-Allow-Methods", web.CONFIG.Cors.AccessMethods)
		c.Header("Access-Control-Expose-Headers", web.CONFIG.Cors.AccessExposeHeaders)
		c.Header("Access-Control-Allow-Credentials", web.CONFIG.Cors.AccessCredentials)

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
