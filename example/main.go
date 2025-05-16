package main

import (
	"log"
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

	engine := s.Engine()
	// add group api v1
	v1 := engine.Group("/api/v1")
	{
		v1.GET("/health", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "OK")
		})
		admin.Group(v1)
	}

	// noitce the static path should not start with /
	// because static path use /*filepath to match all path start with /
	engine.Static("/admin", "./public")
	log.Printf("open: http://%s/admin in your browser\n", s.SystemAddr())

	s.Run()
}
