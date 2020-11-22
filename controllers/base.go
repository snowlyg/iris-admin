package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/easygorm"
)

func GetCommonListSearch(ctx iris.Context) *easygorm.Search {
	offset := libs.ParseInt(ctx.FormValue("page"), 1)
	limit := libs.ParseInt(ctx.FormValue("limit"), 10)
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
