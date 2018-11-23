package controllers

import (
	"IrisYouQiKangApi/models"
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

	models.DB.Model(&models.Orders{}).Count(&oc)
	models.DB.Model(&models.Users{}).Where("is_client = ?", 1).Count(&cc)
	models.DB.Model(&models.Companies{}).Count(&cp)

	hd.Orders = oc
	hd.Clients = cc
	hd.Companies = cp

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(models.ApiJson{Status: true, Data: hd, Msg: "操作成功"})
}
