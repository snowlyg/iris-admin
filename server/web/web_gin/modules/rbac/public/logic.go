package public

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/web/web_gin"
	"github.com/snowlyg/iris-admin/server/web/web_gin/modules/rbac/admin"
	"github.com/snowlyg/iris-admin/server/zap_server"
	multi "github.com/snowlyg/multi/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNameOrPassword = errors.New("用户名或密码错误")
	ErrCaptcha            = errors.New("验证码错误")
)

// GetAccessToken 登录
func GetAccessToken(req *LoginRequest) (*LoginResponse, error) {
	if !store.Verify(req.CaptchaId, req.Captcha, true) && web_gin.CONFIG.System.Level != "test" {
		return nil, ErrCaptcha
	}
	admin, err := admin.FindPasswordByUserName(database.Instance(), req.Username)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		zap_server.ZAPLOG.Error("用户名或密码错误", zap.String("密码:", req.Password), zap.String("hash:", admin.Password), zap.String("bcrypt.CompareHashAndPassword()", err.Error()))
		return nil, ErrUserNameOrPassword
	}

	claims := &multi.CustomClaims{
		ID:            strconv.FormatUint(uint64(admin.Id), 10),
		Username:      req.Username,
		AuthorityId:   "",
		AuthorityType: req.AuthorityType,
		LoginType:     multi.LoginTypeWeb,
		AuthType:      multi.AuthPwd,
		CreationDate:  time.Now().Local().Unix(),
		ExpiresIn:     multi.RedisSessionTimeoutWeb.Milliseconds(),
	}
	token, _, err := multi.AuthDriver.GenerateToken(claims)
	if err != nil {
		return nil, err
	}
	loginResponse := &LoginResponse{
		Data: map[string]interface{}{
			"id": admin.Id,
		},
		Token: token,
	}
	return loginResponse, nil
}

// // GetDeviceAccessToken 登录
// func GetDeviceAccessToken(req *LoginDeviceRequest) (*LoginResponse, error) {

// 	tenancy, err := GetTenancyByUUID(req.UUID)
// 	if err != nil {
// 		return nil, fmt.Errorf("find tenancy %w", err)
// 	}
// 	if tenancy.Status == g.StatusFalse {
// 		return nil, fmt.Errorf("商户已被冻结")
// 	}
// 	cuserId, err := CreateCUserFromDevice(req, tenancy.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	claims := &multi.CustomClaims{
// 		ID:            strconv.FormatUint(uint64(cuserId), 10), // 患者 id
// 		Username:      req.HospitalNo,                          // 用户名使用住院号
// 		TenancyId:     tenancy.ID,
// 		TenancyName:   tenancy.Name,
// 		AuthorityId:   strconv.FormatUint(uint64(g.DeviceAuthorityId), 10),
// 		AuthorityType: multi.GeneralAuthority,
// 		LoginType:     multi.LoginTypeDevice,
// 		AuthType:      multi.AuthPwd,
// 		CreationDate:  time.Now().Local().Unix(),
// 		ExpiresIn:     multi.RedisSessionTimeoutWeb.Milliseconds(),
// 	}
// 	token, _, err := multi.AuthDriver.GenerateToken(claims)
// 	if err != nil {
// 		return nil, err
// 	}
// 	loginResponse := &LoginResponse{
// 		Data: map[string]interface{}{
// 			"id":         cuserId,
// 			"hospitalNo": req.HospitalNo,
// 		},
// 		Token: token,
// 	}
// 	return loginResponse, nil
// }

// DelToken 删除token
func DelToken(token string) error {
	err := multi.AuthDriver.DelUserTokenCache(token)
	if err != nil {
		zap_server.ZAPLOG.Error("del token", zap.Any("err", err))
		return fmt.Errorf("del token %w", err)
	}
	return nil
}

// CleanToken 清空 token
func CleanToken(authorityType int, userId string) error {
	err := multi.AuthDriver.CleanUserTokenCache(authorityType, userId)
	if err != nil {
		zap_server.ZAPLOG.Error("clean token", zap.Any("err", err))
		return fmt.Errorf("clean token %w", err)
	}
	return nil
}

// // GetMiniCode
// func GetMiniCode(miniCode *MiniCode) (*MiniCodeResponse, error) {
// 	mini, err := getMiniProgram()
// 	if err != nil {
// 		return nil, err
// 	}
// 	a := mini.GetAuth()
// 	result, err := a.Code2Session(miniCode.Code)
// 	if err != nil {
// 		return nil, fmt.Errorf("find tenancy %w", err)
// 	}
// 	rmc := &MiniCodeResponse{
// 		SessionKey: result.SessionKey,
// 		OpenId:     result.OpenID,
// 		UnionId:    result.UnionID,
// 	}
// 	return rmc, nil
// }

// // GetMiniAccessToken
// func GetMiniAccessToken(miniLogin *MiniLogin) (*LoginResponse, error) {
// 	tenancy, err := GetTenancyByUUID(miniLogin.UUID)
// 	if err != nil {
// 		return nil, fmt.Errorf("find tenancy %w", err)
// 	}
// 	if tenancy.Status == g.StatusFalse {
// 		return nil, fmt.Errorf("商户已被冻结")
// 	}

// 	mini, err := getMiniProgram()
// 	if err != nil {
// 		return nil, err
// 	}
// 	enc := mini.GetEncryptor()
// 	plainData, err := enc.Decrypt(miniLogin.SessionKey, miniLogin.EncryptedData, miniLogin.Iv)
// 	if err != nil {
// 		return nil, err
// 	}
// 	baseGeneralInfo := model.BaseGeneralInfo{
// 		OpenId:    miniLogin.OpenId,
// 		UnionId:   miniLogin.UnionId,
// 		NickName:  plainData.NickName,
// 		City:      plainData.City,
// 		Province:  plainData.Province,
// 		Country:   plainData.Country,
// 		AvatarUrl: plainData.AvatarURL,
// 		Phone:     plainData.PhoneNumber,
// 		Sex:       plainData.Gender, //1男性,2女性,0未知
// 	}

// 	cuserId, err := CreateCUserFromMini(baseGeneralInfo, tenancy.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	claims := &multi.CustomClaims{
// 		ID:            strconv.FormatUint(uint64(cuserId), 10), // 患者 id
// 		Username:      baseGeneralInfo.NickName,                // 用户微信昵称
// 		TenancyId:     tenancy.ID,
// 		TenancyName:   tenancy.Name,
// 		AuthorityId:   strconv.FormatUint(uint64(g.DeviceAuthorityId), 10),
// 		AuthorityType: multi.GeneralAuthority,
// 		LoginType:     multi.LoginTypeDevice,
// 		AuthType:      multi.AuthPwd,
// 		CreationDate:  time.Now().Local().Unix(),
// 		ExpiresIn:     multi.RedisSessionTimeoutWeb.Milliseconds(),
// 	}
// 	token, _, err := multi.AuthDriver.GenerateToken(claims)
// 	if err != nil {
// 		return nil, err
// 	}
// 	loginResponse := &LoginResponse{
// 		Data: map[string]interface{}{
// 			"nickname": baseGeneralInfo.NickName,
// 		},
// 		Token: token,
// 	}
// 	return loginResponse, nil
// }

// // AuthPhoneNumber
// func AuthPhoneNumber(phoneNumber *PhoneNumber, tenancyId uint) error {
// 	mini, err := getMiniProgram()
// 	if err != nil {
// 		return err
// 	}
// 	enc := mini.GetEncryptor()
// 	plainData, err := enc.Decrypt(phoneNumber.SessionKey, phoneNumber.EncryptedData, phoneNumber.Iv)
// 	if err != nil {
// 		return err
// 	}
// 	baseGeneralInfo := model.BaseGeneralInfo{
// 		OpenId:  phoneNumber.OpenId,
// 		UnionId: phoneNumber.UnionId,
// 		Phone:   plainData.PhoneNumber,
// 	}

// 	_, err = CreateCUserFromMini(baseGeneralInfo, tenancyId)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
