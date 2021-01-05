package auth

import (
	"errors"
	"github.com/bwmarrin/snowflake"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/snowlyg/blog/application/libs"
	"github.com/snowlyg/blog/application/libs/logging"
	"strconv"
	"time"
)

const (
	ZxwSessionTokenPrefix        = "ZST:"
	ZxwSessionBindUserPrefix     = "ZSBU:"
	ZxwSessionUserPrefix         = "ZSU:"
	ZxwSessionUserMaxTokenPrefix = "ZXWUserMaxToken"
)

var (
	ErrTokenInvalid               = errors.New("token is invalid")
	ZxwSessionUserMaxTokenDefault = 10
)

const (
	NoneScope uint64 = iota
	AdminScope
)

const (
	NoAuth int = iota
	AuthPwd
	AuthCode
	AuthThirdParty
)

const (
	LoginTypeWeb int = iota
	LoginTypeApp
	LoginTypeWx
	LoginTypeAlipay
	LoginApplet
)

var (
	RedisSessionTimeoutWeb    = 30 * time.Minute
	RedisSessionTimeoutApp    = 24 * time.Hour
	RedisSessionTimeoutWx     = 5 * 52 * 168 * time.Hour
	RedisSessionTimeoutApplet = 7 * 24 * time.Hour
)

type SessionV2 struct {
	UserId       string `json:"user_id" redis:"user_id"`
	LoginType    int    `json:"login_type" redis:"login_type"`
	AuthType     int    `json:"auth_type" redis:"auth_type"`
	CreationDate int64  `json:"creation_data" redis:"creation_data"`
	ExpiresIn    int    `json:"expires_in" redis:"expires_in"`
	Scope        uint64 `json:"scope" redis:"scope"`
}

// Authentication  认证
type Authentication interface {
	ToCache(token string, id uint64) error
	SyncUserTokenCache(token string) error
	DelUserTokenCache(token string) error
	UserTokenExpired(token string) error
	UpdateUserTokenCacheExpire(token string) error
	GetSessionV2(token string) (*SessionV2, error)
	GetAuthId(token string) (uint, error)
	IsUserTokenOver(token string) bool
	CleanUserTokenCache(token string) error
	Close()
}

var authDriver Authentication

// NewAuthDriver 认证驱动
// redis 需要设置redis
// local 使用本地内存
func NewAuthDriver() Authentication {
	if authDriver != nil {
		return authDriver
	}

	switch libs.Config.Cache.Driver {
	case "redis":
		return NewRedisAuth()
	case "local":
		return NewLocalAuth()
	default:
		return NewLocalAuth()
	}
}

// Login 登录
func Login(id uint64) (string, error) {
	authDriver := NewAuthDriver()
	defer authDriver.Close()
	if authDriver.IsUserTokenOver(strconv.FormatUint(id, 10)) {
		return "", errors.New("以达到同时登录设备上限")
	}
	// 使用分布唯一算法
	node, err := snowflake.NewNode(1)
	if err != nil {
		return "", err
	}
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * time.Duration(1)).Unix(),
		"iat": node.Generate(),
	})
	tokenString, err := token.SignedString([]byte("HS2JDFKhu7Y1av7b"))
	if err != nil {
		logging.ErrorLogger.Errorf("signed string token err", err)
		return "", err
	}
	if err := authDriver.ToCache(tokenString, id); err != nil {
		logging.ErrorLogger.Errorf("to cache user token err", err)
		return "", err
	}
	if err := authDriver.SyncUserTokenCache(tokenString); err != nil {
		logging.ErrorLogger.Errorf("sync user token err", err)
		return "", err
	}
	return tokenString, nil
}

// Logout 退出
func Logout(token string) error {
	authDriver := NewAuthDriver()
	defer authDriver.Close()
	if err := authDriver.DelUserTokenCache(token); err != nil {
		logging.ErrorLogger.Errorf("del user token err", err)
		return err
	}
	return nil
}

// Expire 更新
func Expire(token string) error {
	authDriver := NewAuthDriver()
	defer authDriver.Close()
	if err := authDriver.UpdateUserTokenCacheExpire(token); err != nil {
		logging.ErrorLogger.Errorf("update user token err", err)
		return err
	}
	return nil
}

// Check
func Check(token string) (*SessionV2, error) {
	authDriver := NewAuthDriver()
	defer authDriver.Close()
	rsv2, err := authDriver.GetSessionV2(token)
	if err != nil {
		logging.ErrorLogger.Errorf("check user token err", err)
		return nil, err
	}
	return rsv2, nil
}

// Clear 清除
func Clear(token string) error {
	authDriver := NewAuthDriver()
	defer authDriver.Close()
	err := authDriver.CleanUserTokenCache(token)
	if err != nil {
		logging.ErrorLogger.Errorf("check user token err", err)
		return err
	}
	return nil
}

// AuthId
func AuthId(token string) (uint, error) {
	authDriver := NewAuthDriver()
	defer authDriver.Close()
	id, err := authDriver.GetAuthId(token)
	if err != nil {
		logging.ErrorLogger.Errorf("get user id err", err)
		return 0, err
	}
	return id, err
}
