package perm

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/database/scope"
)

// First 详情
func First(ctx iris.Context) {
	req := &orm.ReqId{}
	if err := req.Request(ctx); err != nil {
		ctx.JSON(orm.Response{Code: orm.ParamErr.Code, Data: nil, Msg: orm.ParamErr.Msg})
		return
	}
	perm := &Response{}
	err := orm.First(database.Instance(), perm, scope.IdScope(req.Id))
	if err != nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: orm.SystemErr.Msg})
		return
	}
	ctx.JSON(orm.Response{Code: orm.NoErr.Code, Data: perm, Msg: orm.NoErr.Msg})
}

// CreatePerm 添加
func CreatePerm(ctx iris.Context) {
	req := &Request{}
	if err := req.Request(ctx); err != nil {
		ctx.JSON(orm.Response{Code: orm.ParamErr.Code, Data: nil, Msg: orm.ParamErr.Msg})
		return
	}
	if !CheckNameAndAct(NameScope(req.Name), ActScope(req.Act)) {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: str.Join("权限[", req.Name, "-", req.Act, "]已存在")})
	}
	perm := &Permission{BasePermission: req.BasePermission}
	id, err := orm.Create(database.Instance(), perm)
	if err != nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(orm.Response{Code: orm.NoErr.Code, Data: iris.Map{"id": id}, Msg: orm.NoErr.Msg})
}

// UpdatePerm 更新
func UpdatePerm(ctx iris.Context) {
	reqId := &orm.ReqId{}
	if err := reqId.Request(ctx); err != nil {
		ctx.JSON(orm.Response{Code: orm.ParamErr.Code, Data: nil, Msg: orm.ParamErr.Msg})
		return
	}
	req := &Request{}
	if err := req.Request(ctx); err != nil {
		ctx.JSON(orm.Response{Code: orm.ParamErr.Code, Data: nil, Msg: orm.ParamErr.Msg})
		return
	}
	if !CheckNameAndAct(NameScope(req.Name), ActScope(req.Act), scope.NeIdScope(reqId.Id)) {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: str.Join("权限[", req.Name, "-", req.Act, "]已存在")})
	}
	perm := &Permission{BasePermission: req.BasePermission}
	err := orm.Update(database.Instance(), reqId.Id, perm)
	if err != nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(orm.Response{Code: orm.NoErr.Code, Data: nil, Msg: orm.NoErr.Msg})
}

// DeletePerm 删除
func DeletePerm(ctx iris.Context) {
	reqId := &orm.ReqId{}
	if err := reqId.Request(ctx); err != nil {
		ctx.JSON(orm.Response{Code: orm.ParamErr.Code, Data: nil, Msg: orm.ParamErr.Msg})
		return
	}
	err := orm.Delete(database.Instance(), reqId.Id, &Permission{})
	if err != nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(orm.Response{Code: orm.NoErr.Code, Data: nil, Msg: orm.NoErr.Msg})
}

// GetAll 分页列表
// - 获取分页参数
// - 请求分页数据
func GetAll(ctx iris.Context) {
	req := &orm.Paginate{}
	if err := req.Request(ctx); err != nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	items := &PageResponse{}
	total, err := orm.Pagination(database.Instance(), items, req.PaginateScope())
	if err != nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	list := iris.Map{"items": items, "total": total, "pageSize": req.PageSize, "page": req.Page}
	ctx.JSON(orm.Response{Code: orm.NoErr.Code, Data: list, Msg: orm.NoErr.Msg})
}
