package auth

import (
	"strconv"
	"time"

	"github.com/snowlyg/iris-admin/modules/user"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/multi"
)

// GetAccessToken 登录
func GetAccessToken(id uint) (string, error) {
	admin, err := user.FindById(database.Instance(), id)
	if err != nil {
		return "", err
	}

	claims := &multi.CustomClaims{
		ID:            strconv.FormatUint(uint64(id), 10),
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
