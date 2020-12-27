package auth

import (
	"errors"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/snowlyg/blog/application/libs/logging"
	"strconv"
	"time"
)

// Login 登录
func Login(auth Authentication, id uint64) (string, error) {
	if auth.IsUserTokenOver(strconv.FormatUint(id, 10)) {
		return "", errors.New("以达到同时登录设备上限")
	}

	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		"iat": time.Now().Unix(),
	})
	tokenString, _ := token.SignedString([]byte("HS2JDFKhu7Y1av7b"))
	if err := auth.ToCache(tokenString, id); err != nil {
		logging.ErrorLogger.Errorf("tocache user token err:%+v", err)
		return "", err
	}

	if err := auth.SyncUserTokenCache(tokenString); err != nil {
		logging.ErrorLogger.Errorf("sync user token err:%+v", err)
		return "", err
	}

	return tokenString, nil

}

// Logout 退出
func Logout(auth Authentication, token string) error {
	if err := auth.DelUserTokenCache(token); err != nil {
		logging.ErrorLogger.Errorf("del user token err:%+v", err)
		return err
	}
	return nil
}

// Expire 更新
func Expire(auth Authentication, token string) error {
	if err := auth.UpdateUserTokenCacheExpire(token); err != nil {
		logging.ErrorLogger.Errorf("update user token err:%+v", err)
		return err
	}
	return nil
}

// Check 更新
func Check(auth Authentication, token string) (*SessionV2, error) {
	rsv2, err := auth.GetSessionV2(token)
	if err != nil {
		logging.ErrorLogger.Errorf("check user token err:%+v", err)
		return nil, err
	}
	return rsv2, nil
}
