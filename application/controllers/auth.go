package controllers

import (
	"fmt"
	"strings"

	"github.com/snowlyg/iris-admin/service/dao"
	"github.com/snowlyg/iris-admin/service/dao/duser"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/jameskeane/bcrypt"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/application/libs"
	"github.com/snowlyg/iris-admin/application/libs/easygorm"
	"github.com/snowlyg/iris-admin/application/libs/logging"
	"github.com/snowlyg/iris-admin/application/libs/response"
	"github.com/snowlyg/iris-admin/application/models"
)

type LoginRe struct {
	Username string `json:"username" validate:"required,gte=2,lte=50" comment:"用户名"`
	Password string `json:"password" validate:"required,gte=6,lte=30"  comment:"密码"`
}

type User struct {
	Id       uint
	Username string
	Password string
}

type Token struct {
	AccessToken string
}

func Login(ctx iris.Context) {
	loginReq := &LoginRe{}
	if err := ctx.ReadJSON(loginReq); err != nil {
		logging.ErrorLogger.Errorf("login read request json err ", err)
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, response.SystemErr.Msg))
		return
	}

	logging.DebugLogger.Debugf("login user ", loginReq)

	validErr := libs.Validate.Struct(*loginReq)
	errs := libs.ValidRequest(validErr)
	if len(errs) > 0 {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, strings.Join(errs, ";")))
		return
	}

	user := User{}
	err := easygorm.GetEasyGormDb().Model(models.User{}).Where("username = ?", loginReq.Username).First(&user).Error
	if err != nil {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}

	logging.DebugLogger.Debugf("user", user)

	if user.Id == 0 {
		ctx.JSON(response.NewResponse(response.DataEmptyErr.Code, nil, fmt.Sprintf("用户 %s 不存在", user.Username)))
		return
	}

	if ok := bcrypt.Match(loginReq.Password, user.Password); !ok {
		ctx.JSON(response.NewResponse(response.AuthErr.Code, nil, "用户名或密码错误"))
		return
	}

	var token string
	token, err = duser.Login(uint64(user.Id))
	if err != nil {
		ctx.JSON(response.NewResponse(response.AuthErr.Code, nil, err.Error()))
		return
	}

	record := models.Oplog{
		Ip:     ctx.RemoteAddr(),
		Method: ctx.Method(),
		Path:   ctx.Path(),
		Agent:  ctx.Request().UserAgent(),
		Status: ctx.GetStatusCode(),
		Body:   "登录后台",
		UserID: user.Id,
	}

	if err := dao.CreateOplog(record); err != nil {
		logging.ErrorLogger.Errorf("create operation record error:", err)
	}

	logging.DebugLogger.Debugf("user token %s", token)

	ctx.JSON(response.NewResponse(response.NoErr.Code, &Token{AccessToken: token}, response.NoErr.Msg))
}

func Logout(ctx iris.Context) {
	jwt, ok := ctx.Values().Get("jwt").(*jwt.Token)
	if !ok {
		_, _ = ctx.JSON(response.NewResponse(response.AuthErr.Code, nil, response.AuthErr.Msg))
		ctx.StopExecution()
		return
	}
	err := duser.Logout(jwt.Raw)
	if err != nil {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}
	ctx.JSON(response.NewResponse(response.NoErr.Code, nil, response.NoErr.Msg))
}

func Expire(ctx iris.Context) {
	jwt, ok := ctx.Values().Get("jwt").(*jwt.Token)
	if !ok {
		_, _ = ctx.JSON(response.NewResponse(response.AuthErr.Code, nil, response.AuthErr.Msg))
		ctx.StopExecution()
		return
	}
	if err := duser.Expire(jwt.Raw); err != nil {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}
	ctx.JSON(response.NewResponse(response.NoErr.Code, nil, response.NoErr.Msg))
}

func Clear(ctx iris.Context) {
	jwt, ok := ctx.Values().Get("jwt").(*jwt.Token)
	if !ok {
		_, _ = ctx.JSON(response.NewResponse(response.AuthErr.Code, nil, response.AuthErr.Msg))
		ctx.StopExecution()
		return
	}
	if err := duser.Clear(jwt.Raw); err != nil {
		ctx.JSON(response.NewResponse(response.SystemErr.Code, nil, err.Error()))
		return
	}
	ctx.JSON(response.NewResponse(response.NoErr.Code, nil, response.NoErr.Msg))
}
