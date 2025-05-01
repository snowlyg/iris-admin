package admin

import (
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func run(address string, router *gin.Engine) serve {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 10 * time.Millisecond
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20
	return s
}
