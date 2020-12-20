package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/jameskeane/bcrypt"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/blog/libs"
	"github.com/snowlyg/blog/libs/logging"
	"github.com/snowlyg/blog/models"
	"github.com/snowlyg/blog/service/auth"
	"github.com/snowlyg/blog/validates"
	"github.com/snowlyg/easygorm"
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

	if ok := bcrypt.Match(loginReq.Password, user.Password); !ok {
		ctx.JSON(libs.ApiResource(libs.ReqErr.Code, nil, "用户名或密码错误"))
		return
	}

	authDriver := auth.NewAuthDriver()
	defer authDriver.Close()
	if tokenString, err := auth.Login(authDriver, uint64(user.ID)); err != nil {
		ctx.JSON(libs.ApiResource(libs.DBErr.Code, nil, err.Error()))
		return
	} else {
		ctx.JSON(libs.ApiResource(libs.NoErr.Code, &models.Token{Token: tokenString}, libs.NoErr.Msg))
		return
	}
}

func UserLogout(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	value := ctx.Values().Get("jwt").(*jwt.Token)
	authDriver := auth.NewAuthDriver()
	defer authDriver.Close()
	err := auth.Logout(authDriver, value.Raw)
	if err != nil {
		ctx.JSON(libs.ApiResource(libs.TokenCacheErr.Code, nil, err.Error()))
		return
	}

	ctx.JSON(libs.ApiResource(libs.NoErr.Code, nil, libs.NoErr.Msg))
	return
}

func UserExpire(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	value := ctx.Values().Get("jwt").(*jwt.Token)
	authDriver := auth.NewAuthDriver()
	defer authDriver.Close()
	if err := auth.Expire(authDriver, value.Raw); err != nil {
		ctx.JSON(libs.ApiResource(libs.TokenCacheErr.Code, nil, libs.TokenCacheErr.Msg))
		return
	}
	ctx.JSON(libs.ApiResource(libs.NoErr.Code, nil, ""))
	return
}
