package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/jameskeane/bcrypt"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/libs/logging"
	"github.com/snowlyg/blog/models"
	"github.com/snowlyg/blog/validates"
	"github.com/snowlyg/easygorm"
	"strconv"
)

func UserLogin(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	loginReq := new(validates.LoginRequest)
	if err := ctx.ReadJSON(loginReq); err != nil {
		logging.Err.Errorf("login read request json err:%+v", err)
		ctx.JSON(libs.ApiResource(libs.ParamReadErr.Code, nil, libs.ParamReadErr.Msg))
		return
	}
	err := validates.Validate.Struct(*loginReq)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				logging.Err.Errorf("login validate request err:%+v", e)
				ctx.JSON(libs.ApiResource(libs.ReqErr.Code, nil, e))
				return
			}
		}
	}
	s := &easygorm.Search{
		Fields: []*easygorm.Field{
			{
				Key:       "username",
				Condition: "=",
				Value:     loginReq.Username,
			},
		},
	}
	user := &models.User{}
	user, err = models.GetUser(s)
	if err != nil && !models.IsNotFound(err) {
		logging.Err.Errorf("login get user info err:%+v", err)
		ctx.JSON(libs.ApiResource(libs.DBErr.Code, nil, libs.DBErr.Msg))
		return
	}

	if user.ID == 0 {
		ctx.JSON(libs.ApiResource(libs.EmptyErr.Code, nil, libs.EmptyErr.Msg))
		return
	}

	uid := strconv.FormatUint(uint64(user.ID), 10)
	if models.IsUserTokenOver(uid) {
		ctx.JSON(libs.ApiResource(libs.ReqErr.Code, nil, "以达到同时登录设备上限"))
		return
	}

	if ok := bcrypt.Match(loginReq.Password, user.Password); !ok {
		ctx.JSON(libs.ApiResource(libs.ReqErr.Code, nil, "用户名或密码错误"))
		return
	}
	token := &models.Token{}
	token.Token, err = models.CacheToken(user)
	if err != nil {
		logging.Err.Errorf("login cache token err:%+v", err)
		ctx.JSON(libs.ApiResource(libs.DBErr.Code, nil, libs.DBErr.Msg))
		return
	}

	ctx.JSON(libs.ApiResource(libs.NoErr.Code, token, libs.NoErr.Msg))
	return
}

func UserLogout(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	value := ctx.Values().Get("jwt").(*jwt.Token)
	conn := libs.GetRedisClusterClient()
	defer conn.Close()
	sess, err := models.GetRedisSessionV2(conn, value.Raw)
	if err != nil {
		logging.Err.Errorf("user logout get redis session err:%+v", err)
		ctx.JSON(libs.ApiResource(libs.TokenCacheErr.Code, nil, libs.TokenCacheErr.Msg))
		return
	}
	if sess != nil {
		if err := sess.DelUserTokenCache(conn, value.Raw); err != nil {
			logging.Err.Errorf("user logout del user token err:%+v", err)
			ctx.JSON(libs.ApiResource(libs.TokenCacheErr.Code, nil, libs.TokenCacheErr.Msg))
			return
		}
	}
	ctx.JSON(libs.ApiResource(libs.NoErr.Code, nil, libs.NoErr.Msg))
	return
}

func UserExpire(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	value := ctx.Values().Get("jwt").(*jwt.Token)
	conn := libs.GetRedisClusterClient()
	defer conn.Close()
	sess, err := models.GetRedisSessionV2(conn, value.Raw)
	if err != nil {
		logging.Err.Errorf("user expire get token err:%+v", err)
		ctx.JSON(libs.ApiResource(libs.TokenCacheErr.Code, nil, libs.TokenCacheErr.Msg))
		return
	}
	if sess != nil {
		if err := sess.UpdateUserTokenCacheExpire(conn, value.Raw); err != nil {
			logging.Err.Errorf("user expire update user token err:%+v", err)
			ctx.JSON(libs.ApiResource(libs.TokenCacheErr.Code, nil, libs.TokenCacheErr.Msg))
			return
		}
	}
	ctx.JSON(libs.ApiResource(libs.NoErr.Code, nil, libs.NoErr.Msg))
	return
}
