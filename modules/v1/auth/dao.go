package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/modules/v1/user"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/multi"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// GetAccessToken 登录
func GetAccessToken(req LoginRequest) (string, error) {
	admin, err := user.FindByUserName(database.Instance(), req.Username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		g.ZAPLOG.Error("用户名或密码错误 ", zap.String("密码:", req.Password), zap.String("hash:", admin.Password), zap.String("错误:", err.Error()))
		return "", errors.New("用户名或密码错误")
	}

	claims := &multi.CustomClaims{
		ID:            strconv.FormatUint(uint64(admin.Id), 10),
		Username:      admin.Username,
		AuthorityId:   "",
		AuthorityType: multi.AdminAuthority,
		LoginType:     multi.LoginTypeWeb,
		AuthType:      multi.AuthPwd,
		CreationDate:  time.Now().Local().Unix(),
		ExpiresIn:     multi.RedisSessionTimeoutWeb.Milliseconds(),
	}
	token, _, err := multi.AuthDriver.GenerateToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

