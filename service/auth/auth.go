package auth

import "github.com/snowlyg/iris-admin/application/libs"

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
