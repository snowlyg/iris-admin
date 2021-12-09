package public

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/snowlyg/iris-admin/server/web/web_gin"
	"github.com/snowlyg/iris-admin/server/web/web_gin/response"
	"github.com/snowlyg/iris-admin/server/zap_server"
	multi "github.com/snowlyg/multi/gin"
	"go.uber.org/zap"
)

var store = base64Captcha.DefaultMemStore

// //ClientLogin 后台登录
// func ClientLogin(ctx *gin.Context) {
// 	req := &LoginRequest{}
// 	if errs := req.Request(ctx); errs != nil {
// 		response.UnauthorizedFailWithMessage(errs.Error(), ctx)
// 		return
// 	}
// 	req.AuthorityType = multi.TenancyAuthority
// 	token, err := GetAccessToken(req)
// 	if err != nil {
// 		zap_server.ZAPLOG.Error("登陆失败!", zap.Any("err", err))
// 		response.UnauthorizedFailWithMessage(err.Error(), ctx)
// 	} else {
// 		response.OkWithData(token, ctx)
// 	}
// }

// AdminLogin 后台登录
func AdminLogin(ctx *gin.Context) {
	req := &LoginRequest{}
	if errs := req.Request(ctx); errs != nil {
		response.UnauthorizedFailWithMessage(errs.Error(), ctx)
		return
	}
	req.AuthorityType = multi.AdminAuthority
	token, err := GetAccessToken(req)
	if err != nil {
		zap_server.ZAPLOG.Error("登陆失败!", zap.Any("err", err))
		response.UnauthorizedFailWithMessage(err.Error(), ctx)
	} else {
		response.OkWithData(token, ctx)
	}
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

// // LoginDevice 床旁用户登录
// func LoginDevice(ctx *gin.Context) {
// 	req := &LoginDeviceRequest{}
// 	if errs := req.Request(ctx); errs != nil {
// 		response.UnauthorizedFailWithMessage(errs.Error(), ctx)
// 		return
// 	}

// 	if token, err := GetDeviceAccessToken(req); err != nil {
// 		zap_server.ZAPLOG.Error("登录失败!", zap.Any("err", err))
// 		response.UnauthorizedFailWithMessage(err.Error(), ctx)
// 	} else {
// 		response.OkWithData(token, ctx)
// 	}
// }

// // AuthMini 小程序用户授权
// func AuthMini(ctx *gin.Context) {
// 	req := &MiniCode{}
// 	if errs := req.Request(ctx); errs != nil {
// 		response.UnauthorizedFailWithMessage(errs.Error(), ctx)
// 		return
// 	}

// 	if sessionKey, err := GetMiniCode(req); err != nil {
// 		zap_server.ZAPLOG.Error("授权失败!", zap.Any("err", err))
// 		response.UnauthorizedFailWithMessage("授权失败"+err.Error(), ctx)
// 	} else {
// 		response.OkWithDetailed(sessionKey, "授权成功", ctx)
// 	}
// }

// // LoginMini 小程序用户授权
// func LoginMini(ctx *gin.Context) {
// 	req := &MiniLogin{}
// 	if errs := req.Request(ctx); errs != nil {
// 		response.UnauthorizedFailWithMessage(errs.Error(), ctx)
// 		return
// 	}

// 	if loginResponse, err := GetMiniAccessToken(req); err != nil {
// 		zap_server.ZAPLOG.Error("登录失败!", zap.Any("err", err))
// 		response.FailWithMessage("登录失败"+err.Error(), ctx)
// 	} else {
// 		response.OkWithDetailed(loginResponse, "登录成功", ctx)
// 	}
// }

// // GetPhoneNumber 小程序用户登录
// func GetPhoneNumber(ctx *gin.Context) {
// 	req := &PhoneNumber{}
// 	if errs := req.Request(ctx); errs != nil {
// 		response.FailWithMessage(errs.Error(), ctx)
// 		return
// 	}

// 	if err := AuthPhoneNumber(req, multi.GetTenancyId(ctx)); err != nil {
// 		zap_server.ZAPLOG.Error("授权失败!", zap.Any("err", err))
// 		response.FailWithMessage("授权失败"+err.Error(), ctx)
// 	} else {
// 		response.OkWithMessage("授权成功", ctx)
// 	}
// }

// Captcha 生成验证码
func Captcha(ctx *gin.Context) {
	//字符,公式,验证码配置
	// 生成默认数字的driver
	driver := base64Captcha.NewDriverDigit(web_gin.CONFIG.Captcha.ImgHeight, web_gin.CONFIG.Captcha.ImgWidth, web_gin.CONFIG.Captcha.KeyLong, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	if id, b64s, err := cp.Generate(); err != nil {
		zap_server.ZAPLOG.Error("验证码获取失败!", zap.Any("err", err))
		response.FailWithMessage(err.Error(), ctx)
	} else {
		response.OkWithData(gin.H{"captchaId": id, "picPath": b64s}, ctx)
	}
}
