package perm

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/database/scope"
)

// First 详情
func First(ctx iris.Context) {
	req := &g.ReqId{}
	if err := req.Request(ctx); err != nil {
		ctx.JSON(g.Response{Code: g.ParamErr.Code, Data: nil, Msg: g.ParamErr.Msg})
		return
	}
	perm := &Response{}
	err := orm.First(database.Instance(), perm, scope.IdScope(req.Id))
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
	if !CheckNameAndAct(NameScope(req.Name), ActScope(req.Act)) {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: str.Join("权限[", req.Name, "-", req.Act, "]已存在")})
	}
	perm := &Permission{BasePermission: req.BasePermission}
	id, err := orm.Create(database.Instance(), perm)
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
	if !CheckNameAndAct(NameScope(req.Name), ActScope(req.Act), scope.NeIdScope(reqId.Id)) {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: str.Join("权限[", req.Name, "-", req.Act, "]已存在")})
	}
	perm := &Permission{BasePermission: req.BasePermission}
	err := orm.Update(database.Instance(), reqId.Id, perm)
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
	err := orm.Delete(database.Instance(), reqId.Id, &Permission{})
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: nil, Msg: g.NoErr.Msg})
}

// GetAll 分页列表
// - 获取分页参数
// - 请求分页数据
func GetAll(ctx iris.Context) {
	req := &g.Paginate{}
	if err := req.Request(ctx); err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	items := &PageResponse{}
	total, err := orm.Paginate(database.Instance(), items, req.PaginateScope())
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	list := iris.Map{"items": items, "total": total, "pageSize": req.PageSize, "page": req.Page}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: list, Msg: g.NoErr.Msg})
}
