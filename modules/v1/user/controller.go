package user

import (
	"errors"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/database/scope"
	"github.com/snowlyg/iris-admin/server/web/validate"
	"github.com/snowlyg/multi"
	"gorm.io/gorm"
)

// Profile 个人信息
func Profile(ctx iris.Context) {
	item := &Response{}
	err := orm.First(database.Instance(), item, scope.IdScope(multi.GetUserId(ctx)))
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: g.SystemErr.Msg})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: item, Msg: g.NoErr.Msg})
}

// GetUser 详情
func GetUser(ctx iris.Context) {
	req := &g.ReqId{}
	if err := req.Request(ctx); err != nil {
		ctx.JSON(g.Response{Code: g.ParamErr.Code, Data: nil, Msg: g.ParamErr.Msg})
		return
	}
	user := &Response{}
	err := orm.First(database.Instance(), user, scope.IdScope(req.Id))
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: g.SystemErr.Msg})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: user, Msg: g.NoErr.Msg})
}

// CreateUser 添加
func CreateUser(ctx iris.Context) {
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

// UpdateUser 更新
func UpdateUser(ctx iris.Context) {
	reqId := &g.ReqId{}
	if err := reqId.Request(ctx); err != nil {
		ctx.JSON(g.Response{Code: g.ParamErr.Code, Data: nil, Msg: g.ParamErr.Msg})
		return
	}

	if err := IsAdminUser(reqId.Id); err != nil {
		ctx.JSON(g.Response{Code: g.ParamErr.Code, Data: nil, Msg: g.ParamErr.Msg})
	}

	req := &Request{}
	if err := req.Request(ctx); err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	if _, err := FindByUserName(UserNameScope(req.Username), scope.NeIdScope(reqId.Id)); !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
	}

	user := &User{BaseUser: req.BaseUser}
	err := orm.Update(database.Instance(), reqId.Id, user)
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	if err := AddRoleForUser(user); err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
	}

	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: nil, Msg: g.NoErr.Msg})
}

// DeleteUser 删除
func DeleteUser(ctx iris.Context) {
	reqId := &g.ReqId{}
	if err := reqId.Request(ctx); err != nil {
		ctx.JSON(g.Response{Code: g.ParamErr.Code, Data: nil, Msg: g.ParamErr.Msg})
		return
	}
	err := orm.Delete(database.Instance(), reqId.Id, &User{})
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
	// 查询用户角色
	getRoles(database.Instance(), *items...)
	list := iris.Map{"items": items, "total": total, "pageSize": req.PageSize, "page": req.Page}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: list, Msg: g.NoErr.Msg})
}

// Logout 退出
func Logout(ctx iris.Context) {
	token := multi.GetVerifiedToken(ctx)
	if token == nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: "授权凭证为空"})
		return
	}
	err := DelToken(string(token))
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: nil, Msg: g.NoErr.Msg})
}

// Clear 清空 token
func Clear(ctx iris.Context) {
	token := multi.GetVerifiedToken(ctx)
	if token == nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: "授权凭证为空"})
		return
	}
	if err := CleanToken(multi.AdminAuthority, string(token)); err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: nil, Msg: g.NoErr.Msg})
}

// ChangeAvatar 修改头像
func ChangeAvatar(ctx iris.Context) {
	avatar := &Avatar{}
	if err := ctx.ReadJSON(avatar); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}
	err := UpdateAvatar(database.Instance(), multi.GetUserId(ctx), avatar.Avatar)
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: nil, Msg: g.NoErr.Msg})
}
