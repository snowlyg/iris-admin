package controllers

import (
	"IrisYouQiKangApi/models"
	"IrisYouQiKangApi/system"
	"github.com/kataras/iris"
	"net/http"
)

type HomeData struct {
	Orders    int `json:"orders"`
	Clients   int `json:"clients"`
	Companies int `json:"companies"`
}

func GetHomeData(ctx iris.Context) {
	oc, cc, cp := 0, 0, 0
	hd := new(HomeData)

	system.DB.Model(&models.Orders{}).Count(&oc)
	system.DB.Model(&models.Users{}).Where("is_client = ?", 1).Count(&cc)
	system.DB.Model(&models.Companies{}).Count(&cp)

	hd.Orders = oc
	hd.Clients = cc
	hd.Companies = cp

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(apiResource(true, hd, "操作成功"))
}
