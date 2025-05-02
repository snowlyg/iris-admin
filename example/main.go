package main

import (
	"github.com/gin-gonic/gin"
	admin "github.com/snowlyg/iris-admin"
)

func main() {
	s, err := admin.NewServe()
	if err != nil {
		panic(err.Error())
	}
	// change default config
	s.Config().System.Addr = "127.0.0.1:8080"
	s.Config().System.GinMode = gin.DebugMode
	if err := s.Config().Recover(); err != nil {
		panic(err.Error())
	}

	s.Run()
}
