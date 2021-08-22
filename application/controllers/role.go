package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/application/libs"
	"github.com/snowlyg/iris-admin/application/libs/logging"
	"github.com/snowlyg/iris-admin/application/libs/response"
	"github.com/snowlyg/iris-admin/service/dao"
	"github.com/snowlyg/iris-admin/service/dao/drole"
)

func GetRole(ctx iris.Context) {
	info := drole.RoleResponse{}
	err := dao.First(&info, ctx)
	if err != nil {
		logging.ErrorLogger.Errorf("get role get err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}
	ctx.JSON(response.NewResponse(response.NoErr.Code, info, response.NoErr.Msg))
}

func CreateRole(ctx iris.Context) {
	roleReq := &drole.RoleReq{}
	if err := ctx.ReadJSON(roleReq); err != nil {
		logging.ErrorLogger.Errorf("create role read json err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, response.SystemErr.Msg))
		return
	}

	validErr := libs.Validate.Struct(*roleReq)
	errs := libs.ValidRequest(validErr)
	if len(errs) > 0 {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, strings.Join(errs, ";")))
		return
	}
	err := dao.Create(&drole.RoleResponse{}, ctx, map[string]interface{}{
		"Name":        roleReq.Name,
		"DisplayName": roleReq.DisplayName,
		"Description": roleReq.Description,
		"CreatedAt":   time.Now(),
		"UpdatedAt":   time.Now(),
	})
	if err != nil {
		logging.ErrorLogger.Errorf("create role get err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}

	ctx.JSON(response.NewResponse(response.NoErr.Code, roleReq, response.NoErr.Msg))
}

func UpdateRole(ctx iris.Context) {
	roleReq := &drole.RoleReq{}
	if err := ctx.ReadJSON(roleReq); err != nil {
		logging.ErrorLogger.Errorf("create role read json err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, response.SystemErr.Msg))
		return
	}

	validErr := libs.Validate.Struct(*roleReq)
	errs := libs.ValidRequest(validErr)
	if len(errs) > 0 {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, strings.Join(errs, ";")))
		return
	}

	err := dao.Update(&drole.RoleResponse{}, ctx, map[string]interface{}{
		"Name":        roleReq.Name,
		"DisplayName": roleReq.DisplayName,
		"Description": roleReq.Description,
		"UpdatedAt":   time.Now(),
	})
	if err != nil {
		logging.ErrorLogger.Errorf("update role read json err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}
	ctx.JSON(response.NewResponse(response.NoErr.Code, nil, response.NoErr.Msg))
}

func DeleteRole(ctx iris.Context) {
	err := dao.Delete(&drole.RoleResponse{}, ctx)
	if err != nil {
		logging.ErrorLogger.Errorf("del role read json err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}

	ctx.JSON(response.NewResponse(response.NoErr.Code, nil, response.NoErr.Msg))
	return
}

func GetAllRoles(ctx iris.Context) {
	name := ctx.FormValue("name")
	page, _ := strconv.Atoi(ctx.FormValue("page"))
	pageSize, _ := strconv.Atoi(ctx.FormValue("pageSize"))
	orderBy := ctx.FormValue("orderBy")
	sort := ctx.FormValue("sort")

	list, err := dao.All(&drole.RoleResponse{}, ctx, name, sort, orderBy, page, pageSize)
	if err != nil {
		logging.ErrorLogger.Errorf("get all role read json err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}
	ctx.JSON(response.NewResponse(response.NoErr.Code, list, response.NoErr.Msg))
}
