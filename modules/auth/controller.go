package user

import (
	"fmt"
	"strconv"
	"time"

	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/modules/user"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/multi"
	"go.uber.org/zap"
)

// Login 登录
func Login(id uint) (string, error) {
	admin, err := user.FindById(database.Instance(), id)
	if err != nil {
		return "", err
	}

	claims := &multi.CustomClaims{
		ID:            strconv.FormatUint(uint64(id), 10),
		Username:      admin.Username,
		AuthorityId:   admin.Roles[],
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

