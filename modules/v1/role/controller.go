package role

import (
	"errors"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/database/scope"
	"gorm.io/gorm"
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

// CreateRole 添加
func CreateRole(ctx iris.Context) {
	req := &Request{}
	if err := req.Request(ctx); err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	id, err := Create(req)
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
	}

	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: iris.Map{"id": id}, Msg: g.NoErr.Msg})
}

// UpdateRole 更新
func UpdateRole(ctx iris.Context) {
	reqId := &g.ReqId{}
	if err := reqId.Request(ctx); err != nil {
		ctx.JSON(g.Response{Code: g.ParamErr.Code, Data: nil, Msg: g.ParamErr.Msg})
		return
	}

	err := IsAdminRole(reqId.Id)
	if err != nil {
		ctx.JSON(g.Response{Code: g.ParamErr.Code, Data: nil, Msg: g.ParamErr.Msg})
	}

	req := &Request{}
	if err := req.Request(ctx); err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	if _, err := FindByName(NameScope(req.Name), scope.NeIdScope(reqId.Id)); !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: "角色名称已经被使用"})
		return
	}

	role := &Role{BaseRole: req.BaseRole}
	err = orm.Update(database.Instance(), reqId.Id, role)
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	err = AddPermForRole(reqId.Id, req.Perms)
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
	}

	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: nil, Msg: g.NoErr.Msg})
}

// DeleteRole 删除
func DeleteRole(ctx iris.Context) {
	reqId := &g.ReqId{}
	if err := reqId.Request(ctx); err != nil {
		ctx.JSON(g.Response{Code: g.ParamErr.Code, Data: nil, Msg: g.ParamErr.Msg})
		return
	}

	err := IsAdminRole(reqId.Id)
	if err != nil {
		ctx.JSON(g.Response{Code: g.ParamErr.Code, Data: nil, Msg: g.ParamErr.Msg})
	}

	err = orm.Delete(database.Instance(), reqId.Id, &Role{})
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: nil, Msg: g.NoErr.Msg})
}

// GetAll 分页列表
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
