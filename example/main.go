package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	admin "github.com/snowlyg/iris-admin"
	"github.com/snowlyg/iris-admin/conf"
)

func main() {
	c := conf.NewConf()
	// change default config
	if err := c.Recover(); err != nil {
		panic(err.Error())
	}
	s, err := admin.NewServe(c)
	if err != nil {
		panic(err.Error())
	}

	// add group api v1
	v1 := s.Engine().Group("/api/v1")
	{
		v1.GET("/health", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "OK")
		})
	}

	// noitce the static path should not start with /
	// because static path use /*filepath to match all path start with /
	s.Engine().Static("/admin", "./public")

	s.Run()
}
