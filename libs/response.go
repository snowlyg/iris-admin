package libs

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/easygorm"
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

type ErrorMsg struct {
	Code int64
	Msg  string
}

var (
	NoErr         = ErrorMsg{200, "请求成功"}
	ReqErr        = ErrorMsg{400, "请求错误"}
	ParamReadErr  = ErrorMsg{5001, "参数解析错误"}
	DBErr         = ErrorMsg{5003, "数据库处理错误"}
	EmptyErr      = ErrorMsg{5004, "数据为空"}
	TokenCacheErr = ErrorMsg{5004, "TOKEN CACHE 错误"}
)

func GetCommonListSearch(ctx iris.Context) *easygorm.Search {
	offset := ParseInt(ctx.FormValue("page"), 1)
	limit := ParseInt(ctx.FormValue("limit"), 10)
	orderBy := ctx.FormValue("orderBy")
	sort := ctx.FormValue("sort")
	field := ctx.FormValue("field")
	relation := ctx.FormValue("relation")

	return &easygorm.Search{
		Sort:      sort,
		Offset:    offset,
		Limit:     limit,
		OrderBy:   orderBy,
		Relations: easygorm.GetRelations(relation, nil),
		Selects:   easygorm.GetSelects(field),
	}
}

func GetCommonSearch(ctx iris.Context) *easygorm.Search {
	orderBy := ctx.FormValue("orderBy")
	sort := ctx.FormValue("sort")
	field := ctx.FormValue("field")
	relation := ctx.FormValue("relation")

	return &easygorm.Search{
		Sort:      sort,
		OrderBy:   orderBy,
		Relations: easygorm.GetRelations(relation, nil),
		Selects:   easygorm.GetSelects(field),
	}
}
