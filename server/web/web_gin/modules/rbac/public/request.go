package public

import (
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"go.uber.org/zap"
)

// LoginRequest 登录请求字段
type LoginRequest struct {
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password" binding:"required"`
	Captcha       string `json:"captcha" binding:"dev-required"`
	CaptchaId     string `json:"captchaId" binding:"dev-required"`
	AuthorityType int    `json:"authorityType" `
}

func (req *LoginRequest) Request(ctx *gin.Context) error {
	if err := ctx.ShouldBindJSON(req); err != nil {
		zap_server.ZAPLOG.Error("参数验证失败", zap.String("ReadParams()", err.Error()))
		return orm.ErrParamValidate
	}
	return nil
}

type LoginDeviceRequest struct {
	UUID       string `json:"uuid" binding:"required"`
	Phone      string `json:"phone" form:"phone" binding:"required"`
	Name       string `json:"name" form:"name" binding:"required"`
	Sex        int    `json:"sex" form:"sex"`
	Age        int    `json:"age" form:"age"`
	LocName    string `json:"locName" form:"locName"  binding:"required"`
	BedNum     string `json:"bedNum" form:"bedNum" binding:"required"`
	HospitalNo string `json:"hospitalNo" form:"hospitalNo" binding:"required"`
	Disease    string `json:"disease" form:"disease"  binding:"required"`
}

func (req *LoginDeviceRequest) Request(ctx *gin.Context) error {
	if err := ctx.ShouldBindJSON(req); err != nil {
		zap_server.ZAPLOG.Error("参数验证失败", zap.String("ReadParams()", err.Error()))
		return orm.ErrParamValidate
	}
	return nil
}

type MiniCode struct {
	Code string `json:"code" form:"code" binding:"required"`
}

func (req *MiniCode) Request(ctx *gin.Context) error {
	if err := ctx.ShouldBindJSON(req); err != nil {
		zap_server.ZAPLOG.Error("参数验证失败", zap.String("ReadParams()", err.Error()))
		return orm.ErrParamValidate
	}
	return nil
}

type MiniLogin struct {
	UUID          string `json:"uuid"  form:"uuid" binding:"required"`
	SessionKey    string `json:"sessionKey" form:"sessionKey" binding:"required"`
	Iv            string `json:"iv"  form:"iv" binding:"required"`
	EncryptedData string `json:"encryptedData" form:"encryptedData" binding:"required"`
	OpenId        string `json:"openId" form:"openId" binding:"required`
	UnionId       string `json:"unionId" form:"unionId" binding:"required`
}

func (req *MiniLogin) Request(ctx *gin.Context) error {
	if err := ctx.ShouldBindJSON(req); err != nil {
		zap_server.ZAPLOG.Error("参数验证失败", zap.String("ReadParams()", err.Error()))
		return orm.ErrParamValidate
	}
	return nil
}

type PhoneNumber struct {
	Iv            string `json:"iv"  form:"iv" binding:"required"`
	EncryptedData string `json:"encryptedData" form:"encryptedData" binding:"required"`
	SessionKey    string `json:"sessionKey" form:"sessionKey" binding:"required"`
	OpenId        string `json:"openId" form:"openId" binding:"required`
	UnionId       string `json:"unionId" form:"unionId" binding:"required`
}

func (req *PhoneNumber) Request(ctx *gin.Context) error {
	if err := ctx.ShouldBindJSON(req); err != nil {
		zap_server.ZAPLOG.Error("参数验证失败", zap.String("ReadParams()", err.Error()))
		return orm.ErrParamValidate
	}
	return nil
}
