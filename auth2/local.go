package auth2

import (
	"errors"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type tokens []string

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

// GenerateToken
func (la *LocalAuth) GenerateToken(claims *MultiClaims) (string, int64, error) {
	if la.isUserTokenOver(claims.AuthorityType, claims.Id) {
		return "", 0, ErrOverMaxTokenCount
	}
	token, err := GetToken()
	if err != nil {
		return "", 0, err
	}
	err = la.toCache(token, claims)
	if err != nil {
		return "", 0, err
	}
	if err = la.syncUserTokenCache(token); err != nil {
		return "", 0, err
	}

	return token, int64(claims.ExpiresAt), err
}

func (la *LocalAuth) toCache(token string, rcc *MultiClaims) error {
	sKey := GtSessionTokenPrefix + token
	la.Cache.Set(sKey, rcc, GetTokenExpire(rcc.LoginType))
	return nil
}

func (la *LocalAuth) syncUserTokenCache(token string) error {
	rcc, err := la.GetMultiClaims(token)
	if err != nil {
		return err
	}

	userPrefixKey := GetUserPrefixKey(rcc.AuthorityType, rcc.Id)
	ts := tokens{}
	if uTokens, ok := la.Cache.Get(userPrefixKey); ok && uTokens != nil {
		ts = uTokens.(tokens)
	}
	ts = append(ts, token)
	la.Cache.Set(userPrefixKey, ts, cache.NoExpiration)

	la.Cache.Set(GtSessionBindUserPrefix+token, userPrefixKey, GetTokenExpire(rcc.LoginType))
	return nil
}

func (la *LocalAuth) DelUserTokenCache(token string) error {
	rcc, err := la.GetMultiClaims(token)
	if err != nil {
		return err
	}
	if rcc == nil {
		return errors.New("token cache is nil")
	}

	userPrefixKey := GetUserPrefixKey(rcc.AuthorityType, rcc.Id)
	if utokens, ok := la.Cache.Get(userPrefixKey); ok && utokens != nil {
		t := utokens.(tokens)
		for index, u := range t {
			if u == token {
				if len(t) == 1 {
					utokens = nil
				} else {
					utokens = append(t[0:index], t[index:]...)
				}
			}
		}
		la.Cache.Set(userPrefixKey, utokens, cache.NoExpiration)

	}
	err = la.delTokenCache(token)
	if err != nil {
		return err
	}

	return nil
}

// delTokenCache 删除token缓存
func (la *LocalAuth) delTokenCache(token string) error {
	la.Cache.Delete(GtSessionBindUserPrefix + token)
	la.Cache.Delete(GtSessionTokenPrefix + token)
	return nil
}

func (la *LocalAuth) UpdateUserTokenCacheExpire(token string) error {
	rsv2, err := la.GetMultiClaims(token)
	if err != nil {
		return err
	}
	if rsv2 == nil {
		return errors.New("token cache is nil")
	}
	la.Cache.Set(GtSessionBindUserPrefix+token, rsv2, GetTokenExpire(rsv2.LoginType))
	la.Cache.Set(GtSessionTokenPrefix+token, rsv2, GetTokenExpire(rsv2.LoginType))

	return nil
}

func (la *LocalAuth) GetMultiClaims(token string) (*MultiClaims, error) {
	sKey := GtSessionTokenPrefix + token
	if food, found := la.Cache.Get(sKey); !found || food == nil {
		return nil, ErrTokenInvalid
	} else {
		return food.(*MultiClaims), nil
	}
}

// GetTokenByClaims 获取用户信息
func (la *LocalAuth) GetTokenByClaims(cla *MultiClaims) (string, error) {
	userTokens, err := la.getUserTokens(cla.AuthorityType, cla.Id)
	if err != nil {
		return "", err
	}
	clas, err := la.getMultiClaimses(userTokens)
	if err != nil {
		return "", err
	}
	for token, existCla := range clas {
		if cla.AuthType == existCla.AuthType &&
			cla.Id == existCla.Id &&
			cla.AuthorityType == existCla.AuthorityType &&
			cla.AuthorityId == existCla.AuthorityId &&
			cla.LoginType == existCla.LoginType {
			return token, nil
		}
	}
	return "", nil
}

// getUserTokens 获取登录数量
func (la *LocalAuth) getUserTokens(authorityType int, userId string) (tokens, error) {
	if utokens, ok := la.Cache.Get(GetUserPrefixKey(authorityType, userId)); ok && utokens != nil {
		return utokens.(tokens), nil
	}
	return nil, nil
}

// getMultiClaimses 获取用户信息
func (la *LocalAuth) getMultiClaimses(tokens tokens) (map[string]*MultiClaims, error) {
	clas := make(map[string]*MultiClaims, la.getUserTokenMaxCount())
	for _, token := range tokens {
		cla, err := la.GetMultiClaims(token)
		if err != nil {
			continue
		}
		clas[token] = cla
	}

	return clas, nil
}

func (la *LocalAuth) isUserTokenOver(authorityType int, userId string) bool {
	return la.getUserTokenCount(authorityType, userId) >= la.getUserTokenMaxCount()
}

// getUserTokenCount 获取登录数量
func (la *LocalAuth) getUserTokenCount(authorityType int, userId string) int64 {
	return la.checkMaxCount(authorityType, userId)
}

func (la *LocalAuth) checkMaxCount(authorityType int, userId string) int64 {
	utokens, _ := la.getUserTokens(authorityType, userId)
	if utokens == nil {
		return 0
	}
	for index, u := range utokens {
		if _, found := la.Cache.Get(GtSessionTokenPrefix + u); !found {
			if len(utokens) == 1 {
				utokens = nil
			} else {
				utokens = append(utokens[0:index], utokens[index:]...)
			}
		}
	}
	la.Cache.Set(GetUserPrefixKey(authorityType, userId), utokens, cache.NoExpiration)
	return int64(len(utokens))

}

// getUserTokenMaxCount 最大登录限制
func (la *LocalAuth) getUserTokenMaxCount() int64 {
	if count, found := la.Cache.Get(GtSessionUserMaxTokenPrefix); !found {
		return GtSessionUserMaxTokenDefault
	} else {
		return count.(int64)
	}
}

// SetUserTokenMaxCount 最大登录限制
func (la *LocalAuth) SetUserTokenMaxCount(tokenMaxCount int64) error {
	la.Cache.Set(GtSessionUserMaxTokenPrefix, tokenMaxCount, cache.NoExpiration)
	return nil
}

// CleanUserTokenCache 清空token缓存
func (la *LocalAuth) CleanUserTokenCache(authorityType int, userId string) error {
	utokens, _ := la.getUserTokens(authorityType, userId)
	if utokens == nil {
		return nil
	}

	for _, token := range utokens {
		err := la.delTokenCache(token)
		if err != nil {
			continue
		}
	}
	la.Cache.Delete(GetUserPrefixKey(authorityType, userId))

	return nil
}

// IsRole
func (la *LocalAuth) IsRole(token string, authorityType int) (bool, error) {
	rcc, err := la.GetMultiClaims(token)
	if err != nil {
		return false, fmt.Errorf("get User's infomation return error: %w", err)
	}
	return rcc.AuthorityType == authorityType, nil
}

// IsRole
func (la *LocalAuth) IsSuperAdmin(token string) bool {
	rcc, err := la.GetMultiClaims(token)
	if err != nil {
		return false
	}
	return rcc.SuperAdmin
}

// 兼容 redis
func (la *LocalAuth) Close() {}
