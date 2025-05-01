package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func run(address string, router *gin.Engine) serve {
	return &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
