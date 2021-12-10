package auth

import (
	"errors"

	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/web/web_iris/modules/rbac/user"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"github.com/snowlyg/multi"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNameOrPassword = errors.New("用户名或密码错误")
)

// GetAccessToken 登录
func GetAccessToken(req *LoginRequest) (string, error) {
	admin, err := user.FindPasswordByUserName(database.Instance(), req.Username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		zap_server.ZAPLOG.Error("用户名或密码错误", zap.String("密码:", req.Password), zap.String("hash:", admin.Password), zap.String("bcrypt.CompareHashAndPassword()", err.Error()))
		return "", ErrUserNameOrPassword
	}

	claims := multi.New(&multi.Multi{
		Id:            admin.Id,
		Username:      req.Username,
		AuthorityIds:  admin.AuthorityIds,
		AuthorityType: multi.AdminAuthority,
		LoginType:     multi.LoginTypeWeb,
		AuthType:      multi.AuthPwd,
		ExpiresIn:     multi.RedisSessionTimeoutWeb.Milliseconds(),
	})
	token, _, err := multi.AuthDriver.GenerateToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}
