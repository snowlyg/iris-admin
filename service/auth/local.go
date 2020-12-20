package auth

import (
	"errors"
	"github.com/patrickmn/go-cache"
	"github.com/snowlyg/blog/libs/logging"
	"strconv"
	"strings"
	"time"
)

type tokens []string
type skeys []string

var localCache *cache.Cache

type LocalAuth struct {
	Cache *cache.Cache
}

func NewLocalAuth() *LocalAuth {
	if localCache == nil {
		localCache = cache.New(4*time.Hour, 24*time.Minute)
	}
	return &LocalAuth{
		Cache: localCache,
	}
}

func (la *LocalAuth) ToCache(token string, id uint64) error {
	sKey := ZXW_SESSION_TOKEN_PREFIX + token
	rsv2 := &SessionV2{
		UserId:       strconv.FormatUint(id, 10),
		LoginType:    LoginTypeWeb,
		AuthType:     AuthPwd,
		CreationDate: time.Now().Unix(),
		Scope:        GetUserScope("admin"),
	}
	la.Cache.Set(sKey, rsv2, la.getTokenExpire(rsv2))
	return nil
}
func (la *LocalAuth) SyncUserTokenCache(token string) error {
	rsv2, err := la.GetSessionV2(token)
	if err != nil {
		logging.Err.Errorf("SyncUserTokenCache err: %+v\n", err)
		return err
	}

	sKey := ZXW_SESSION_USER_PREFIX + rsv2.UserId
	ts := tokens{}
	if uTokens, uFound := la.Cache.Get(sKey); uFound {
		ts = uTokens.(tokens)
	}
	ts = append(ts, token)

	sKey2 := ZXW_SESSION_BIND_USER_PREFIX + token
	sys := skeys{}
	if keys, found := la.Cache.Get(sKey2); found {
		sys = keys.(skeys)
	}
	sys = append(sys, sKey)

	return nil
}
func (la *LocalAuth) DelUserTokenCache(token string) error {
	rsv2, err := la.GetSessionV2(token)
	if err != nil {
		return err
	}
	if rsv2 == nil {
		return errors.New("token cache is nil")
	}
	sKey := ZXW_SESSION_USER_PREFIX + rsv2.UserId
	exp := la.getTokenExpire(rsv2)
	if utokens, ufound := la.Cache.Get(sKey); ufound {
		t := utokens.(tokens)
		for index, u := range t {
			if u == token {
				utokens = append(t[0:index], t[index:]...)
				la.Cache.Set(sKey, utokens, exp)
			}
		}
	}
	err = la.DelTokenCache(token)
	if err != nil {
		return err
	}

	return nil
}

// DelTokenCache 删除token缓存
func (la *LocalAuth) DelTokenCache(token string) error {
	la.Cache.Delete(ZXW_SESSION_BIND_USER_PREFIX + token)
	la.Cache.Delete(ZXW_SESSION_TOKEN_PREFIX + token)
	return nil
}

func (la *LocalAuth) UserTokenExpired(token string) error {
	rsv2, err := la.GetSessionV2(token)
	if err != nil {
		return err
	}
	if rsv2 == nil {
		return errors.New("token cache is nil")
	}

	exp := la.getTokenExpire(rsv2)
	uKey := ZXW_SESSION_BIND_USER_PREFIX + token
	if sKeys, found := la.Cache.Get(uKey); !found {
		return errors.New("token skey is empty")
	} else {
		for _, v := range sKeys.(skeys) {
			if !strings.Contains(v, ZXW_SESSION_USER_PREFIX) {
				continue
			}
			if utokens, ufound := la.Cache.Get(v); ufound {
				t := utokens.(tokens)
				for index, u := range t {
					if u == token {
						utokens = append(t[0:index], t[index:]...)
						la.Cache.Set(v, utokens, exp)
					}
				}
			}
		}
	}

	la.Cache.Delete(uKey)
	return nil
}

func (la *LocalAuth) UpdateUserTokenCacheExpire(token string) error {
	rsv2, err := la.GetSessionV2(token)
	if err != nil {
		return err
	}
	if rsv2 == nil {
		return errors.New("token cache is nil")
	}
	la.Cache.Set(ZXW_SESSION_TOKEN_PREFIX+token, rsv2, la.getTokenExpire(rsv2))

	return nil
}

// getTokenExpire 过期时间
func (la *LocalAuth) getTokenExpire(rsv2 *SessionV2) time.Duration {
	timeout := RedisSessionTimeoutApp
	if rsv2.LoginType == LoginTypeWeb {
		timeout = RedisSessionTimeoutWeb
	} else if rsv2.LoginType == LoginTypeWx {
		timeout = RedisSessionTimeoutWx
	} else if rsv2.LoginType == LoginTypeAlipay {
		timeout = RedisSessionTimeoutWx
	}
	return timeout
}

func (la *LocalAuth) GetSessionV2(token string) (*SessionV2, error) {
	sKey := ZXW_SESSION_TOKEN_PREFIX + token
	get, _ := la.Cache.Get(sKey)
	logging.Dbug.Infof("GetSessionV2: %+v", get)
	if food, found := la.Cache.Get(sKey); !found {
		return nil, ERR_TOKEN_INVALID
	} else {
		return food.(*SessionV2), nil
	}
}

func (la *LocalAuth) IsUserTokenOver(userId string) bool {
	if la.getUserTokenCount(userId) >= la.getUserTokenMaxCount() {
		return true
	}
	return false
}

// getUserTokenCount 获取登录数量
func (la *LocalAuth) getUserTokenCount(userId string) int {
	if userTokens, found := la.Cache.Get(ZXW_SESSION_USER_PREFIX + userId); !found {
		return 0
	} else {
		return len(userTokens.(tokens))
	}
}

// getUserTokenMaxCount 最大登录限制
func (la *LocalAuth) getUserTokenMaxCount() int {
	if count, found := la.Cache.Get(ZXW_SESSION_USER_MAX_TOKEN_PREFIX); !found {
		return ZXW_SESSION_USER_MAX_TOKEN_DEFAULT
	} else {
		return count.(int)
	}
}

// CleanUserTokenCache 清空token缓存
func (la *LocalAuth) CleanUserTokenCache(token string) error {
	//rsv2,err := ra.GetSessionV2(token)
	//if err != nil {
	//	logging.Err.Errorf("clean user token cache member err: %+v", err)
	//	return err
	//}
	//sKey := ZXW_SESSION_USER_PREFIX + rsv2.UserId
	//var allTokens []string
	//allTokens, err = redis.Strings(ra.Conn.Members(sKey))
	//if err != nil {
	//	logging.Err.Errorf("clean user token cache member err: %+v", err)
	//	return err
	//}
	//_, err = ra.Conn.Del(sKey)
	//if err != nil {
	//	logging.Err.Errorf("clean user token cache del err: %+v", err)
	//	return err
	//}
	//
	//for _, token := range allTokens {
	//	err = ra.DelTokenCache(token)
	//	if err != nil {
	//		return err
	//	}
	//}
	return nil
}

// 兼容 redis
func (la *LocalAuth) Close() {}
