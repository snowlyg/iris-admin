package models

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/snowlyg/blog/libs"
	"strings"
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
	RedisSessionTimeoutWeb    = 30 * time.Minute
	RedisSessionTimeoutApp    = 24 * time.Hour
	RedisSessionTimeoutApplet = 7 * 24 * time.Hour
	RedisSessionTimeoutWx     = 5 * 52 * 168 * time.Hour
)

type RedisSessionV2 struct {
	UserId       string `json:"user_id" redis:"user_id"`
	LoginType    int    `json:"login_type" redis:"login_type"`
	AuthType     int    `json:"auth_type" redis:"auth_type"`
	CreationDate int64  `json:"creation_data" redis:"creation_data"`
	ExpiresIn    int    `json:"expires_in" redis:"expires_in"`
	Scope        uint64 `json:"scope" redis:"scope"`
}

//  GetRedisSessionV2获取 session
func GetRedisSessionV2(conn *libs.RedisCluster, token string) (*RedisSessionV2, error) {
	sKey := ZXW_SESSION_TOKEN_PREFIX + token
	if !conn.Exists(sKey) {
		return nil, ERR_TOKEN_INVALID
	}
	pp := new(RedisSessionV2)
	if err := conn.LoadRedisHashToStruct(sKey, pp); err != nil {
		return nil, err
	}
	return pp, nil
}

// isUserTokenOver 超过登录设备限制
func isUserTokenOver(userId string) bool {
	conn := libs.GetRedisClusterClient()
	defer conn.Close()
	if getUserTokenCount(conn, userId) >= getUserTokenMaxCount(conn) {
		return true
	}
	return false
}

// getUserTokenCount 获取登录数量
func getUserTokenCount(conn *libs.RedisCluster, userId string) int {
	count, err := redis.Int(conn.Scard(ZXW_SESSION_USER_PREFIX + userId))
	if err != nil {
		fmt.Println(fmt.Sprintf("getUserTokenCount error :%+v", err))
		return 0
	}
	return count
}

// getUserTokenMaxCount 最大登录限制
func getUserTokenMaxCount(conn *libs.RedisCluster) int {
	count, err := redis.Int(conn.GetKey(ZXW_SESSION_USER_MAX_TOKEN_PREFIX))
	if err != nil {
		return ZXW_SESSION_USER_MAX_TOKEN_DEFAULT
	}
	return count
}

// UserTokenExpired 过期 token
func UserTokenExpired(token string) {
	conn := libs.GetRedisClusterClient()
	defer conn.Close()

	uKey := ZXW_SESSION_BIND_USER_PREFIX + token
	sKeys, err := redis.Strings(conn.Members(uKey))
	if err != nil {
		fmt.Println(fmt.Sprintf("conn.Members key %s error :%+v", uKey, err))
		return
	}
	for _, v := range sKeys {
		if !strings.Contains(v, ZXW_SESSION_USER_PREFIX) {
			continue
		}
		_, err := conn.Do("SREM", v, token)
		if err != nil {
			fmt.Println(fmt.Sprintf("conn.SREM key %s token %s  error :%+v", v, token, err))
			return
		}
	}
	if _, err := conn.Del(uKey); err != nil {
		fmt.Println(fmt.Sprintf("conn.Del key %s error :%+v", uKey, err))
	}
	return
}

// getUserScope 角色
func getUserScope(userType string) uint64 {
	switch userType {
	case "admin":
		return AdminScope
	}
	return NoneScope
}

// ToCache 缓存 token
func (r *RedisSessionV2) ToCache(conn *libs.RedisCluster, token string) error {
	sKey := ZXW_SESSION_TOKEN_PREFIX + token

	if _, err := conn.HMSet(sKey,
		"user_id", r.UserId,
		"login_type", r.LoginType,
		"auth_type", r.AuthType,
		"creation_data", r.CreationDate,
		"expires_in", r.ExpiresIn,
		"scope", r.Scope,
	); err != nil {
		fmt.Println(fmt.Sprintf("conn.ToCache error :%+v", err))
		return err
	}
	return nil
}

// SyncUserTokenCache 同步 token 到缓存
func (r *RedisSessionV2) SyncUserTokenCache(conn *libs.RedisCluster, token string) error {
	sKey := ZXW_SESSION_USER_PREFIX + r.UserId
	if _, err := conn.Sadd(sKey, token); err != nil {
		fmt.Println(fmt.Sprintf("conn.SyncUserTokenCache1 error :%+v", err))
		return err
	}
	sKey2 := ZXW_SESSION_BIND_USER_PREFIX + token
	_, err := conn.Sadd(sKey2, sKey)
	if err != nil {
		fmt.Println(fmt.Sprintf("conn.SyncUserTokenCache2 error :%+v", err))
		return err
	}
	return nil
}

//UpdateUserTokenCacheExpire 更新过期时间
func (r *RedisSessionV2) UpdateUserTokenCacheExpire(conn *libs.RedisCluster, token string) error {
	if _, err := conn.Expire(ZXW_SESSION_TOKEN_PREFIX+token, int(r.GetTokenExpire().Seconds())); err != nil {
		fmt.Println(fmt.Sprintf("conn.UpdateUserTokenCacheExpire error :%+v", err))
		return err
	}
	return nil
}

// GetTokenExpire 过期时间
func (r *RedisSessionV2) GetTokenExpire() time.Duration {
	timeout := RedisSessionTimeoutApp
	if r.LoginType == LoginTypeWeb {
		timeout = RedisSessionTimeoutWeb
	} else if r.LoginType == LoginTypeWx {
		timeout = RedisSessionTimeoutWx
	} else if r.LoginType == LoginTypeAlipay {
		timeout = RedisSessionTimeoutWx
	}
	return timeout
}

// DelUserTokenCache 删除token缓存
func (r *RedisSessionV2) DelUserTokenCache(conn *libs.RedisCluster, token string) error {
	sKey := ZXW_SESSION_USER_PREFIX + r.UserId
	_, err := conn.Do("SREM", sKey, token)
	if err != nil {
		fmt.Println(fmt.Sprintf("conn.DelUserTokenCache1 error :%+v", err))
		return err
	}
	err = r.DelTokenCache(conn, token)
	if err != nil {
		fmt.Println(fmt.Sprintf("conn.DelUserTokenCache2 error :%+v", err))
		return err
	}

	return nil
}

// DelTokenCache 删除token缓存
func (r *RedisSessionV2) DelTokenCache(conn *libs.RedisCluster, token string) error {
	sKey2 := ZXW_SESSION_BIND_USER_PREFIX + token
	_, err := conn.Del(sKey2)
	if err != nil {
		fmt.Println(fmt.Sprintf("conn.DelUserTokenCache2 error :%+v", err))
		return err
	}

	sKey3 := ZXW_SESSION_TOKEN_PREFIX + token
	_, err = conn.Del(sKey3)
	if err != nil {
		fmt.Println(fmt.Sprintf("conn.DelUserTokenCache3 error :%+v", err))
		return err
	}

	return nil
}

// CleanUserTokenCache 清空token缓存
func (r *RedisSessionV2) CleanUserTokenCache(conn *libs.RedisCluster) error {
	sKey := ZXW_SESSION_USER_PREFIX + r.UserId
	allTokens, err := redis.Strings(conn.Members(sKey))
	if err != nil {
		fmt.Println(fmt.Sprintf("conn.CleanUserTokenCache1 error :%+v", err))
		return err
	}
	_, err = conn.Del(sKey)
	if err != nil {
		fmt.Println(fmt.Sprintf("conn.CleanUserTokenCache2 error :%+v", err))
		return err
	}

	for _, token := range allTokens {
		err = r.DelTokenCache(conn, token)
		if err != nil {
			fmt.Println(fmt.Sprintf("conn.DelUserTokenCache2 error :%+v", err))
			return err
		}
	}
	return nil
}
