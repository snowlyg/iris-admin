package admin

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/database/scope"
	"github.com/snowlyg/iris-admin/server/web/web_gin/response"
	multi "github.com/snowlyg/multi/gin"
	"gorm.io/gorm"
)

// Profile 个人信息
func Profile(ctx *gin.Context) {
	item := &Response{}
	err := orm.First(database.Instance(), item, scope.IdScope(multi.GetUserId(ctx)))
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.OkWithData(item, ctx)
}

// GetAdmin 详情
func GetAdmin(ctx *gin.Context) {
	req := &orm.ReqId{}
	if errs := ctx.ShouldBindJSON(&req); errs != nil {
		response.FailWithMessage(errs.Error(), ctx)
		return
	}
	admin := &Response{}
	err := orm.First(database.Instance(), admin, scope.IdScope(req.Id))
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.OkWithData(admin, ctx)
}

// CreateAdmin添加
func CreateAdmin(ctx *gin.Context) {
	req := &Request{}
	if errs := ctx.ShouldBindJSON(&req); errs != nil {
		response.FailWithMessage(errs.Error(), ctx)
		return
	}

	id, err := Create(req)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
	}
	response.OkWithData(gin.H{"id": id}, ctx)
}

// UpdateAdmin 更新
func UpdateAdmin(ctx *gin.Context) {
	reqId := &orm.ReqId{}
	if errs := ctx.ShouldBindJSON(&reqId); errs != nil {
		response.FailWithMessage(errs.Error(), ctx)
		return
	}

	if err := IsAdminUser(reqId.Id); err != nil {
		response.FailWithMessage(err.Error(), ctx)
	}

	req := &Request{}
	if errs := ctx.ShouldBindJSON(&req); errs != nil {
		response.FailWithMessage(errs.Error(), ctx)
		return
	}

	if _, err := FindByUserName(UserNameScope(req.Username), scope.NeIdScope(reqId.Id)); !errors.Is(err, gorm.ErrRecordNotFound) {
		response.FailWithMessage(err.Error(), ctx)
	}

	admin := &Admin{BaseAdmin: req.BaseAdmin}
	err := orm.Update(database.Instance(), reqId.Id, admin)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	if err := AddRoleForUser(admin); err != nil {
		response.FailWithMessage(err.Error(), ctx)
	}

	response.Ok(ctx)
}

// DeleteAdmin 删除
func DeleteAdmin(ctx *gin.Context) {
	reqId := &orm.ReqId{}
	if errs := ctx.ShouldBindJSON(&reqId); errs != nil {
		response.FailWithMessage(errs.Error(), ctx)
		return
	}
	err := orm.Delete(database.Instance(), &Admin{}, scope.IdScope(reqId.Id))
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Ok(ctx)
}

// GetAll 分页列表
func GetAll(ctx *gin.Context) {
	req := &orm.Paginate{}
	if errs := ctx.ShouldBindJSON(&req); errs != nil {
		response.FailWithMessage(errs.Error(), ctx)
		return
	}

	items := &PageResponse{}
	total, err := orm.Pagination(database.Instance(), items, req.PaginateScope())
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	// 查询用户角色
	getRoles(database.Instance(), items.Item...)
	response.OkWithData(response.PageResult{
		List:     items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, ctx)
}

// Logout 退出
func Logout(ctx *gin.Context) {
	token := multi.GetVerifiedToken(ctx)
	if token == nil {
		response.FailWithMessage("授权凭证为空", ctx)
		return
	}
	err := DelToken(string(token))
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Ok(ctx)
}

// Clear 清空 token
func Clear(ctx *gin.Context) {
	token := multi.GetVerifiedToken(ctx)
	if token == nil {
		response.FailWithMessage("授权凭证为空", ctx)
		return
	}
	if err := CleanToken(multi.AdminAuthority, string(token)); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Ok(ctx)
}

// ChangeAvatar 修改头像
func ChangeAvatar(ctx *gin.Context) {
	avatar := &Avatar{}
	if errs := ctx.ShouldBindJSON(&avatar); errs != nil {
		response.FailWithMessage(errs.Error(), ctx)
		return
	}
	err := UpdateAvatar(database.Instance(), multi.GetUserId(ctx), avatar.HeaderImg)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Ok(ctx)
}
