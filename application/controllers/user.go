package controllers

import (
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/application/libs"
	"github.com/snowlyg/blog/application/libs/easygorm"
	"github.com/snowlyg/blog/application/libs/logging"
	"github.com/snowlyg/blog/application/libs/response"
	"github.com/snowlyg/blog/application/models"
	"github.com/snowlyg/blog/service/auth"
	"github.com/snowlyg/blog/service/dao"
	"github.com/snowlyg/blog/service/dao/duser"
	"strconv"
	"strings"
	"time"
)

func Profile(ctx iris.Context) {
	token := ctx.Values().Get("jwt").(*jwt.Token).Raw
	id, err := auth.AuthId(token)
	if err != nil {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, response.SystemErr.Msg))
		return
	}

	info := duser.UserResponse{}
	err = dao.Find(&info, ctx, id)
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

	token := ctx.Values().Get("jwt").(*jwt.Token).Raw
	id, err := auth.AuthId(token)
	if err != nil {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, response.SystemErr.Msg))
		return
	}

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
	err = easygorm.GetEasyGormDb().Model(&models.User{}).Where("id = ?", id).Update("avatar", avatar.Avatar).Error
	if err != nil {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, response.SystemErr.Msg))
		return
	}
	ctx.JSON(response.NewResponse(response.NoErr.Code, nil, response.NoErr.Msg))
}

func GetUser(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	info := duser.UserResponse{}
	err := dao.Find(&info, ctx, id)
	if err != nil {
		logging.ErrorLogger.Errorf("get user get err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}
	ctx.JSON(response.NewResponse(response.NoErr.Code, info, response.NoErr.Msg))
}

func CreateUser(ctx iris.Context) {
	userReq := &duser.UserReq{}
	if err := ctx.ReadJSON(userReq); err != nil {
		logging.ErrorLogger.Errorf("create user read json err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, response.SystemErr.Msg))
		return
	}

	validErr := libs.Validate.Struct(*userReq)
	errs := libs.ValidRequest(validErr)
	if len(errs) > 0 {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, strings.Join(errs, ";")))
		return
	}

	err := dao.Create(&duser.UserResponse{}, ctx, map[string]interface{}{
		"Name":      userReq.Name,
		"Username":  userReq.Username,
		"Password":  libs.HashPassword(userReq.Password),
		"Intro":     userReq.Intro,
		"Avatar":    userReq.Avatar,
		"CreatedAt": time.Now(),
		"UpdatedAt": time.Now(),
	})
	if err != nil {
		logging.ErrorLogger.Errorf("create user get err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}

	ctx.JSON(response.NewResponse(response.NoErr.Code, userReq, response.NoErr.Msg))
	return

}

func UpdateUser(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	userReq := &duser.UserReq{}
	if err := ctx.ReadJSON(userReq); err != nil {
		logging.ErrorLogger.Errorf("create user read json err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, response.SystemErr.Msg))
		return
	}

	validErr := libs.Validate.Struct(*userReq)
	errs := libs.ValidRequest(validErr)
	if len(errs) > 0 {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, strings.Join(errs, ";")))
		return
	}

	err := dao.Update(&duser.UserResponse{}, ctx, id, map[string]interface{}{
		"Name":      userReq.Name,
		"Password":  libs.HashPassword(userReq.Password),
		"Intro":     userReq.Intro,
		"Avatar":    userReq.Avatar,
		"UpdatedAt": time.Now(),
	})
	if err != nil {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}
	ctx.JSON(response.NewResponse(response.NoErr.Code, userReq, response.NoErr.Msg))
	return
}

func DeleteUser(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	err := dao.Delete(&duser.UserResponse{}, ctx, id)
	if err != nil {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}

	ctx.JSON(response.NewResponse(response.NoErr.Code, nil, response.NoErr.Msg))
	return
}

// GetUsers
func GetUsers(ctx iris.Context) {
	name := ctx.FormValue("name")
	page, _ := strconv.Atoi(ctx.FormValue("page"))
	pageSize, _ := strconv.Atoi(ctx.FormValue("pageSize"))
	orderBy := ctx.FormValue("orderBy")
	sort := ctx.FormValue("sort")

	list, err := dao.All(&duser.UserResponse{}, ctx, name, sort, orderBy, page, pageSize)
	if err != nil {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}
	ctx.JSON(response.NewResponse(response.NoErr.Code, list, response.NoErr.Msg))
	return
}

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
