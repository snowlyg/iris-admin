package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/models"
)

type Response struct {
	Code int64       `json:"code"`
	Msg  interface{} `json:"message"`
	Data interface{} `json:"data"`
}

type Lists struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
}

func ApiResource(code int64, objects interface{}, msg string) (r *Response) {
	r = &Response{Code: code, Data: objects, Msg: msg}
	return
}

func GetCommonListSearch(ctx iris.Context) *models.Search {
	offset := libs.ParseInt(ctx.FormValue("page"), 1)
	limit := libs.ParseInt(ctx.FormValue("limit"), 20)
	orderBy := ctx.FormValue("orderBy")
	sort := ctx.FormValue("sort")
	return &models.Search{
		Sort:    sort,
		Offset:  offset,
		Limit:   limit,
		OrderBy: orderBy,
	}
}
