package main

import (
	"github.com/kataras/iris"
	"net/http"
)

type HomeData struct {
	Orders    int `json:"orders"`
	Clients   int `json:"clients"`
	Companies int `json:"companies"`
}

func CGetHomeData(ctx iris.Context) {
	hd := new(HomeData)
	hd.Orders = MGetOrderCounts()
	hd.Clients = MGetClientCounts()
	hd.Companies = MGetCompanyCounts()

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(apiResource(true, hd, "操作成功"))
}
