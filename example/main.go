package main

import (
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
	s.Engine().Static("/", "./public")
	s.Run()
}
