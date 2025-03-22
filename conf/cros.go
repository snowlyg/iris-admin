package conf

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CorsConf struct {
	AccessOrigin        string `mapstructure:"access-origin" json:"burst" access-origin:"access-origin"`
	AccessHeaders       string `mapstructure:"access-headers" json:"access-headers" yaml:"access-headers"`
	AccessMethods       string `mapstructure:"access-methods" json:"access-methods" yaml:"access-methods"`
	AccessExposeHeaders string `mapstructure:"access-expose-headers" json:"access-expose-headers" yaml:"access-expose-headers"`
	AccessCredentials   string `mapstructure:"access-credentials" json:"access-credentials" yaml:"access-credentials"`
}

// Cors
func (corsConf *CorsConf) Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", corsConf.AccessOrigin)
		c.Header("Access-Control-Allow-Headers", corsConf.AccessHeaders)
		c.Header("Access-Control-Allow-Methods", corsConf.AccessMethods)
		c.Header("Access-Control-Expose-Headers", corsConf.AccessExposeHeaders)
		c.Header("Access-Control-Allow-Credentials", corsConf.AccessCredentials)
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
