package web_gin

import (
	"fmt"
	"syscall"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func initServer(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.BeforeBegin = func(add string) {
		fmt.Printf("Actual pid is %d\n", syscall.Getpid())
		// save it somehow
	}
	s.ReadHeaderTimeout = 10 * time.Millisecond
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20
	return s
}
