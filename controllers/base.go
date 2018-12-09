package controllers

import (
	"IrisApiProject/tools"
	"github.com/kataras/iris"
)

type ApiJson struct {
	Status bool        `json:"status"`
	Msg    interface{} `json:"msg"`
	Data   interface{} `json:"data"`
}

var (
	t *tools.Tools
)

func init() {
	t = tools.New()
}

func NotFound(ctx iris.Context) {
	ctx.JSON(apiResource(false, nil, "404 Not Found"))
}
func InternalServerError(ctx iris.Context) {
	ctx.WriteString("Oups something went wrong, try again")
}

func apiResource(status bool, data interface{}, msg string) (apijson *ApiJson) {
	apijson = &ApiJson{Status: status, Data: data, Msg: msg}
	return
}
