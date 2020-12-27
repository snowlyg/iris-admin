package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/application/libs"
	"github.com/snowlyg/blog/application/libs/easygorm"
	"github.com/snowlyg/blog/application/libs/logging"
	"github.com/snowlyg/blog/application/libs/response"
	"github.com/snowlyg/blog/application/models"
	"github.com/snowlyg/blog/service/auth"
	"strings"
)

type Info struct {
	Id       uint
	Name     string
	Username string
	Intro    string
	Avatar   string
}

func Profile(ctx iris.Context) {
	sess := ctx.Values().Get("sess").(*auth.SessionV2)
	id := uint(libs.ParseInt(sess.UserId, 10))
	info := Info{Id: id}
	err := easygorm.EasyGorm.DB.Model(&models.User{}).Find(&info).Error
	if err != nil {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, response.SystemErr.Msg))
		return
	}
	ctx.JSON(response.NewResponse(response.NoErr.Code, info, response.NoErr.Msg))
}

type Avatar struct {
	Avatar string
}

func ChangeAvatar(ctx iris.Context) {
	sess := ctx.Values().Get("sess").(*auth.SessionV2)
	id := uint(libs.ParseInt(sess.UserId, 10))

	avatar := &Avatar{}
	if err := ctx.ReadJSON(avatar); err != nil {
		logging.ErrorLogger.Errorf("change avatar read json error ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, response.SystemErr.Msg))
		return
	}

	validErr := libs.Validate.Struct(*avatar)
	errs := libs.ValidRequest(validErr)
	if len(errs) > 0 {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, strings.Join(errs, ";")))
		return
	}
	err := easygorm.EasyGorm.DB.Model(&models.User{}).Where("id = ?", id).Update("avatar", avatar.Avatar).Error
	if err != nil {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, response.SystemErr.Msg))
		return
	}
	ctx.JSON(response.NewResponse(response.NoErr.Code, nil, response.NoErr.Msg))
}

//
//func GetUser(ctx iris.Context) {
//	ctx.StatusCode(iris.StatusOK)
//	id, _ := ctx.Params().GetUint("id")
//	s := &easygorm.Search{
//		Fields: []*easygorm.Field{
//			{
//				Key:       "id",
//				Condition: "=",
//				Value:     id,
//			},
//		},
//	}
//	user, err := models.GetUser(s)
//	if err != nil {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, err.Error()))
//		return
//	}
//	_, _ = ctx.JSON(response.NewResponse(200, user, "操作成功"))
//}
//
//func CreateUser(ctx iris.Context) {
//	ctx.StatusCode(iris.StatusOK)
//	user := new(models.User)
//	if err := ctx.ReadJSON(user); err != nil {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, err.Error()))
//		return
//	}
//
//	err := validates.Validate.Struct(*user)
//	if err != nil {
//		errs := err.(validator.ValidationErrors)
//		for _, e := range errs.Translate(validates.ValidateTrans) {
//			if len(e) > 0 {
//				_, _ = ctx.JSON(response.NewResponse(400, nil, e))
//				return
//			}
//		}
//	}
//
//	err = user.CreateUser()
//	if err != nil {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, err.Error()))
//		return
//	}
//
//	if user.ID == 0 {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, "操作失败"))
//		return
//	}
//	_, _ = ctx.JSON(response.NewResponse(200, user, "操作成功"))
//	return
//
//}
//
//func UpdateUser(ctx iris.Context) {
//	ctx.StatusCode(iris.StatusOK)
//	user := new(models.User)
//
//	if err := ctx.ReadJSON(user); err != nil {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, err.Error()))
//	}
//
//	err := validates.Validate.Struct(*user)
//	if err != nil {
//		errs := err.(validator.ValidationErrors)
//		for _, e := range errs.Translate(validates.ValidateTrans) {
//			if len(e) > 0 {
//				_, _ = ctx.JSON(response.NewResponse(400, nil, e))
//				return
//			}
//		}
//	}
//
//	id, _ := ctx.Params().GetUint("id")
//	if user.Username == "username" {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, "不能编辑管理员"))
//		return
//	}
//
//	err = models.UpdateUserById(id, user)
//	if err != nil {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, err.Error()))
//		return
//	}
//	_, _ = ctx.JSON(response.NewResponse(200, user, "操作成功"))
//}
//
//func DeleteUser(ctx iris.Context) {
//	ctx.StatusCode(iris.StatusOK)
//	id, _ := ctx.Params().GetUint("id")
//
//	err := models.DeleteUser(id)
//	if err != nil {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, err.Error()))
//		return
//	}
//	_, _ = ctx.JSON(response.NewResponse(200, nil, "删除成功"))
//}

/**
* @api {get} /users 获取所有的账号
* @apiName 获取所有的账号
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 获取所有的账号
* @apiSampleRequest /users
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
//func GetAllUsers(ctx iris.Context) {
//	ctx.StatusCode(iris.StatusOK)
//	name := ctx.FormValue("name")
//
//	s := libs.GetCommonListSearch(ctx)
//	s.Fields = append(s.Fields, easygorm.GetField("name", name))
//	users, count, err := models.GetAllUsers(s)
//	if err != nil {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, err.Error()))
//		return
//	}
//
//	_, _ = ctx.JSON(response.NewResponse(200, map[string]interface{}{"items": users, "total": count, "limit": s.Limit}, "操作成功"))
//
//}
//
//func usersTransform(users []*models.User) []*transformer.User {
//	var us []*transformer.User
//	for _, user := range users {
//		u := userTransform(user)
//		us = append(us, u)
//	}
//	return us
//}
//
//func userTransform(user *models.User) *transformer.User {
//	u := &transformer.User{}
//	g := gf.NewTransform(u, user, time.RFC3339)
//	_ = g.Transformer()
//
//	roleIds := easygorm.GetRolesForUser(user.ID)
//	var ris []int
//	for _, roleId := range roleIds {
//		ri, _ := strconv.Atoi(roleId)
//		ris = append(ris, ri)
//	}
//	s := &easygorm.Search{
//		Fields: []*easygorm.Field{
//			{
//				Key:       "id",
//				Condition: "IN",
//				Value:     ris,
//			},
//		},
//	}
//	roles, _, err := models.GetAllRoles(s)
//	if err == nil {
//		u.Roles = rolesTransform(roles)
//	}
//	return u
//}
