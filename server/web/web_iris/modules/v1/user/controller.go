package user

import (
	"errors"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/database/scope"
	"github.com/snowlyg/iris-admin/server/web/web_iris/validate"
	"github.com/snowlyg/multi"
	"gorm.io/gorm"
)

// Profile 个人信息
func Profile(ctx iris.Context) {
	item := &Response{}
	err := orm.First(database.Instance(), item, scope.IdScope(multi.GetUserId(ctx)))
	if err != nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: orm.SystemErr.Msg})
		return
	}
	ctx.JSON(orm.Response{Code: orm.NoErr.Code, Data: item, Msg: orm.NoErr.Msg})
}

// GetUser 详情
func GetUser(ctx iris.Context) {
	req := &orm.ReqId{}
	if err := req.Request(ctx); err != nil {
		ctx.JSON(orm.Response{Code: orm.ParamErr.Code, Data: nil, Msg: orm.ParamErr.Msg})
		return
	}
	user := &Response{}
	err := orm.First(database.Instance(), user, scope.IdScope(req.Id))
	if err != nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: orm.SystemErr.Msg})
		return
	}
	ctx.JSON(orm.Response{Code: orm.NoErr.Code, Data: user, Msg: orm.NoErr.Msg})
}

// CreateUser 添加
func CreateUser(ctx iris.Context) {
	req := &Request{}
	if err := req.Request(ctx); err != nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	id, err := Create(req)
	if err != nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: err.Error()})
	}

	ctx.JSON(orm.Response{Code: orm.NoErr.Code, Data: iris.Map{"id": id}, Msg: orm.NoErr.Msg})
}

// UpdateUser 更新
func UpdateUser(ctx iris.Context) {
	reqId := &orm.ReqId{}
	if err := reqId.Request(ctx); err != nil {
		ctx.JSON(orm.Response{Code: orm.ParamErr.Code, Data: nil, Msg: orm.ParamErr.Msg})
		return
	}

	if err := IsAdminUser(reqId.Id); err != nil {
		ctx.JSON(orm.Response{Code: orm.ParamErr.Code, Data: nil, Msg: orm.ParamErr.Msg})
	}

	req := &Request{}
	if err := req.Request(ctx); err != nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	if _, err := FindByUserName(UserNameScope(req.Username), scope.NeIdScope(reqId.Id)); !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: err.Error()})
	}

	user := &User{BaseUser: req.BaseUser}
	err := orm.Update(database.Instance(), reqId.Id, user)
	if err != nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	if err := AddRoleForUser(user); err != nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: err.Error()})
	}

	ctx.JSON(orm.Response{Code: orm.NoErr.Code, Data: nil, Msg: orm.NoErr.Msg})
}

// DeleteUser 删除
func DeleteUser(ctx iris.Context) {
	reqId := &orm.ReqId{}
	if err := reqId.Request(ctx); err != nil {
		ctx.JSON(orm.Response{Code: orm.ParamErr.Code, Data: nil, Msg: orm.ParamErr.Msg})
		return
	}
	err := orm.Delete(database.Instance(), reqId.Id, &User{})
	if err != nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(orm.Response{Code: orm.NoErr.Code, Data: nil, Msg: orm.NoErr.Msg})
}

// GetAll 分页列表
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
	// 查询用户角色
	getRoles(database.Instance(), *items...)
	list := iris.Map{"items": items, "total": total, "pageSize": req.PageSize, "page": req.Page}
	ctx.JSON(orm.Response{Code: orm.NoErr.Code, Data: list, Msg: orm.NoErr.Msg})
}

// Logout 退出
func Logout(ctx iris.Context) {
	token := multi.GetVerifiedToken(ctx)
	if token == nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: "授权凭证为空"})
		return
	}
	err := DelToken(string(token))
	if err != nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(orm.Response{Code: orm.NoErr.Code, Data: nil, Msg: orm.NoErr.Msg})
}

// Clear 清空 token
func Clear(ctx iris.Context) {
	token := multi.GetVerifiedToken(ctx)
	if token == nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: "授权凭证为空"})
		return
	}
	if err := CleanToken(multi.AdminAuthority, string(token)); err != nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(orm.Response{Code: orm.NoErr.Code, Data: nil, Msg: orm.NoErr.Msg})
}

// ChangeAvatar 修改头像
func ChangeAvatar(ctx iris.Context) {
	avatar := &Avatar{}
	if err := ctx.ReadJSON(avatar); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}
	err := UpdateAvatar(database.Instance(), multi.GetUserId(ctx), avatar.Avatar)
	if err != nil {
		ctx.JSON(orm.Response{Code: orm.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(orm.Response{Code: orm.NoErr.Code, Data: nil, Msg: orm.NoErr.Msg})
}
