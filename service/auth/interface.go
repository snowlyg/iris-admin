package auth

import (
	"errors"
	"time"

	"github.com/snowlyg/blog/application/libs"
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

type Session struct {
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
	GetSessionV2(token string) (*Session, error)
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
