package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/application/libs"
	"github.com/snowlyg/blog/application/libs/logging"
	"github.com/snowlyg/blog/application/libs/response"
	"github.com/snowlyg/blog/service/dao"
	"github.com/snowlyg/blog/service/dao/dperm"
)

func GetPermission(ctx iris.Context) {
	info := dperm.PermResponse{}
	err := dao.First(&info, ctx)
	if err != nil {
		logging.ErrorLogger.Errorf("get perm get err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}
	ctx.JSON(response.NewResponse(response.NoErr.Code, info, response.NoErr.Msg))
}

func CreatePermission(ctx iris.Context) {
	permReq := &dperm.PermReq{}
	if err := ctx.ReadJSON(permReq); err != nil {
		logging.ErrorLogger.Errorf("create perm read json err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, response.SystemErr.Msg))
		return
	}

	validErr := libs.Validate.Struct(*permReq)
	errs := libs.ValidRequest(validErr)
	if len(errs) > 0 {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, strings.Join(errs, ";")))
		return
	}
	err := dao.Create(&dperm.PermResponse{}, ctx, map[string]interface{}{
		"Name":        permReq.Name,
		"DisplayName": permReq.DisplayName,
		"Description": permReq.Description,
		"Act":         permReq.Act,
		"CreatedAt":   time.Now(),
		"UpdatedAt":   time.Now(),
	})
	if err != nil {
		logging.ErrorLogger.Errorf("create perm get err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}

	ctx.JSON(response.NewResponse(response.NoErr.Code, permReq, response.NoErr.Msg))
}

func UpdatePermission(ctx iris.Context) {
	permReq := &dperm.PermReq{}
	if err := ctx.ReadJSON(permReq); err != nil {
		logging.ErrorLogger.Errorf("create perm read json err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, response.SystemErr.Msg))
		return
	}

	validErr := libs.Validate.Struct(*permReq)
	errs := libs.ValidRequest(validErr)
	if len(errs) > 0 {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, strings.Join(errs, ";")))
		return
	}

	err := dao.Update(&dperm.PermResponse{}, ctx, map[string]interface{}{
		"Name":        permReq.Name,
		"DisplayName": permReq.DisplayName,
		"Description": permReq.Description,
		"UpdatedAt":   time.Now(),
	})
	if err != nil {
		logging.ErrorLogger.Errorf("update perm read json err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}
	ctx.JSON(response.NewResponse(response.NoErr.Code, nil, response.NoErr.Msg))
}

func DeletePermission(ctx iris.Context) {
	err := dao.Delete(&dperm.PermResponse{}, ctx)
	if err != nil {
		logging.ErrorLogger.Errorf("delete perm read json err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}

	ctx.JSON(response.NewResponse(response.NoErr.Code, nil, response.NoErr.Msg))
}

func GetAllPermissions(ctx iris.Context) {
	name := ctx.FormValue("name")
	page, _ := strconv.Atoi(ctx.FormValue("page"))
	pageSize, _ := strconv.Atoi(ctx.FormValue("pageSize"))
	orderBy := ctx.FormValue("orderBy")
	sort := ctx.FormValue("sort")

	list, err := dao.All(&dperm.PermResponse{}, ctx, name, sort, orderBy, page, pageSize)
	if err != nil {
		logging.ErrorLogger.Errorf("get all perm read json err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}
	ctx.JSON(response.NewResponse(response.NoErr.Code, list, response.NoErr.Msg))
}
