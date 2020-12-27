package auth

import (
	"errors"
	"github.com/snowlyg/blog/application/libs"
	"time"
)

const (
	ZXW_SESSION_TOKEN_PREFIX          = "ZST:"
	ZXW_SESSION_BIND_USER_PREFIX      = "ZSBU:"
	ZXW_SESSION_USER_PREFIX           = "ZSU:"
	ZXW_SESSION_USER_MAX_TOKEN_PREFIX = "ZXWUserMaxToken"
)

var (
	ERR_TOKEN_INVALID                  = errors.New("token is invalid!")
	ZXW_SESSION_USER_MAX_TOKEN_DEFAULT = 10
)

const (
	NoneScope uint64 = iota
	AdminScope
)

const (
	NonoAuth int = iota
	AuthPwd
	AuthCode
	AuthThirdparty
)

const (
	LoginTypeWeb int = iota
	LoginTypeApp
	LoginTypeWx
	LoginTypeAlipay
	LoginApplet
)

var (
	RedisSessionTimeoutWeb = 30 * time.Minute
	RedisSessionTimeoutApp = 24 * time.Hour
	RedisSessionTimeoutWx  = 5 * 52 * 168 * time.Hour
	//RedisSessionTimeoutApplet = 7 * 24 * time.Hour
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
	IsUserTokenOver(token string) bool
	CleanUserTokenCache(token string) error
	Close()
}

// NewAuthDriver 认证驱动
// redis 需要设置redis
// local 使用本地内存
func NewAuthDriver() Authentication {
	switch libs.Config.Cache.Driver {
	case "redis":
		return NewRedisAuth()
	case "local":
		return NewLocalAuth()
	default:
		return NewLocalAuth()
	}
}
