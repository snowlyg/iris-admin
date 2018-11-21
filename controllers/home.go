package controllers

import (
	"github.com/kataras/iris"
)

type HomeData struct {
	Orders    string `json:"orders"`
	Clients   string `json:"clients"`
	Companies string `json:"companies"`
}

func GetHomeData(ctx iris.Context) {
}
