package perm

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database/scope"
)

// GetPerm 详情
func GetPerm(ctx iris.Context) {
	req := &g.ReqId{}
	if err := req.Request(ctx); err != nil {
		ctx.JSON(g.Response{Code: g.ParamErr.Code, Data: nil, Msg: g.ParamErr.Msg})
		return
	}
	perm := GetResponse()
	err := perm.First(scope.IdScope(req.Id))
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: g.SystemErr.Msg})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: perm, Msg: g.NoErr.Msg})
}

// CreatePerm 添加
func CreatePerm(ctx iris.Context) {
	req := &Request{}
	if err := req.Request(ctx); err != nil {
		ctx.JSON(g.Response{Code: g.ParamErr.Code, Data: nil, Msg: g.ParamErr.Msg})
		return
	}
	id, err := Create(req)
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: iris.Map{"id": id}, Msg: g.NoErr.Msg})
}

// UpdatePerm 更新
func UpdatePerm(ctx iris.Context) {
	reqId := &g.ReqId{}
	if err := reqId.Request(ctx); err != nil {
		ctx.JSON(g.Response{Code: g.ParamErr.Code, Data: nil, Msg: g.ParamErr.Msg})
		return
	}
	req := &Request{}
	if err := req.Request(ctx); err != nil {
		ctx.JSON(g.Response{Code: g.ParamErr.Code, Data: nil, Msg: g.ParamErr.Msg})
		return
	}
	err := Update(reqId.Id, req)
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: nil, Msg: g.NoErr.Msg})
}

// DeletePerm 删除
func DeletePerm(ctx iris.Context) {
	reqId := &g.ReqId{}
	if err := reqId.Request(ctx); err != nil {
		ctx.JSON(g.Response{Code: g.ParamErr.Code, Data: nil, Msg: g.ParamErr.Msg})
		return
	}
	err := DeleteById(reqId.Id)
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: nil, Msg: g.NoErr.Msg})
}

// GetAllPerms 分页列表
// - 获取分页参数
// - 请求分页数据
func GetAllPerms(ctx iris.Context) {
	req := &g.Paginate{}
	if err := req.Request(ctx); err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	perms := PageResponse{}
	total, err := perms.Paginate(req.PaginateScope())
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	list := iris.Map{"items": perms, "total": total, "pageSize": req.PageSize, "page": req.Page}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: list, Msg: g.NoErr.Msg})
}
