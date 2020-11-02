package controllers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/models"
	"github.com/snowlyg/blog/validates"
)

/**
* @api {get} /admin/configs/:key 根据id获取权限信息
* @apiName 根据id获取权限信息
* @apiGroup Configs
* @apiVersion 1.0.0
* @apiDescription 根据id获取权限信息
* @apiSampleRequest /admin/configs/:key
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func GetConfig(ctx iris.Context) {

	ctx.StatusCode(iris.StatusOK)
	key := ctx.Params().GetString("key")
	s := &models.Search{
		Fields: []*models.Filed{
			{
				Key:       "name",
				Condition: "=",
				Value:     key,
			},
		},
	}
	config, err := models.GetConfig(s)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, fmt.Sprintf("Error get %s config: %s", key, err.Error())))
		return
	}
	_, _ = ctx.JSON(libs.ApiResource(200, config, "操作成功"))
}

/**
* @api {post} /admin/configs/ 新建权限
* @apiName 新建权限
* @apiGroup Configs
* @apiVersion 1.0.0
* @apiDescription 新建权限
* @apiSampleRequest /admin/configs/
* @apiParam {string} name 权限名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func CreateConfig(ctx iris.Context) {

	ctx.StatusCode(iris.StatusOK)
	config := new(models.Config)
	if err := ctx.ReadJSON(config); err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}
	err := validates.Validate.Struct(*config)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				_, _ = ctx.JSON(libs.ApiResource(400, nil, e))
				return
			}
		}
	}

	err = config.CreateConfig()
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, fmt.Sprintf("Error create prem: %s", err.Error())))
		return
	}

	if config.ID == 0 {
		_, _ = ctx.JSON(libs.ApiResource(400, config, "操作失败"))
		return
	}
	_, _ = ctx.JSON(libs.ApiResource(200, config, "操作成功"))

}

/**
* @api {post} /admin/configs/:id/update 更新权限
* @apiName 更新权限
* @apiGroup Configs
* @apiVersion 1.0.0
* @apiDescription 更新权限
* @apiSampleRequest /admin/configs/:id/update
* @apiParam {string} name 权限名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func UpdateConfig(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	config := new(models.Config)

	if err := ctx.ReadJSON(config); err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}
	err := validates.Validate.Struct(*config)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				_, _ = ctx.JSON(libs.ApiResource(400, nil, e))
				return
			}
		}
	}

	id, _ := ctx.Params().GetUint("id")
	err = models.UpdateConfig(id, config)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, fmt.Sprintf("Error create prem: %s", err.Error())))
		return
	}

	if config.ID == 0 {
		_, _ = ctx.JSON(libs.ApiResource(400, config, "操作失败"))
		return
	}
	_, _ = ctx.JSON(libs.ApiResource(200, config, "操作成功"))

}

/**
* @api {delete} /admin/configs/:id/delete 删除权限
* @apiName 删除权限
* @apiGroup Configs
* @apiVersion 1.0.0
* @apiDescription 删除权限
* @apiSampleRequest /admin/configs/:id/delete
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func DeleteConfig(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	id, _ := ctx.Params().GetUint("id")
	err := models.DeleteConfig(id)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}
	_, _ = ctx.JSON(libs.ApiResource(200, nil, "删除成功"))
}

/**
* @api {get} /configs 获取所有的权限
* @apiName 获取所有的权限
* @apiGroup Configs
* @apiVersion 1.0.0
* @apiDescription 获取所有的权限
* @apiSampleRequest /configs
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func GetAllConfigs(ctx iris.Context) {

	ctx.StatusCode(iris.StatusOK)
	s := GetCommonListSearch(ctx)

	configs, err := models.GetAllConfigs(s)
	if err != nil {
		_, _ = ctx.JSON(libs.ApiResource(400, nil, err.Error()))
		return
	}
	_, _ = ctx.JSON(libs.ApiResource(200, configs, "操作成功"))
}
