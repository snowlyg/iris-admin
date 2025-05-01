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

func NewLocal() *LocalAuth {
	if localCache == nil {
		localCache = cache.New(4*time.Hour, 24*time.Minute)
	}
	return &LocalAuth{
		Cache: localCache,
	}
}

// Generate
func (la *LocalAuth) Generate(claims *Claims) (string, int64, error) {

	if la.isUserTokenOver(claims.roleType(), claims.Id) {
		return "", 0, fmt.Errorf("local: is user token over fail:%w", ErrOverLimit)
	}
	token, err := getToken()
	if err != nil {
		return "", 0, fmt.Errorf("local: get token fail:%w", err)
	}
	la.toCache(token, claims)
	if e := la.syncUserCache(token); e != nil {
		return "", 0, fmt.Errorf("local: sync user token fail:%w", e)
	}
	return token, int64(claims.ExpiresAt), nil
}

func (la *LocalAuth) toCache(token string, rcc *Claims) error {
	sKey := TokenPrefix + token
	la.Cache.Set(sKey, rcc, getExpire(rcc.loginType()))
	return nil
}

func (la *LocalAuth) syncUserCache(token string) error {
	rcc, err := la.GetClaims(token)
	if err != nil {
		return err
	}
	userPrefixKey := getPrefixKey(rcc.roleType(), rcc.Id)
	ts := tokens{}
	if uTokens, ok := la.Cache.Get(userPrefixKey); ok && uTokens != nil {
		ts = uTokens.(tokens)
	}
	ts = append(ts, token)
	la.Cache.Set(userPrefixKey, ts, cache.NoExpiration)
	la.Cache.Set(BindUserPrefix+token, userPrefixKey, getExpire(rcc.loginType()))
	return nil
}

func (la *LocalAuth) DelCache(token string) error {
	rcc, err := la.GetClaims(token)
	if err != nil {
		return err
	}
	userPrefixKey := getPrefixKey(rcc.roleType(), rcc.Id)
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
	la.delTokenCache(token)
	return nil
}

// delTokenCache
func (la *LocalAuth) delTokenCache(token string) error {
	la.Cache.Delete(BindUserPrefix + token)
	la.Cache.Delete(TokenPrefix + token)
	return nil
}

func (la *LocalAuth) UpdateCacheExpire(token string) error {
	rsv2, err := la.GetClaims(token)
	if err != nil {
		return err
	}
	if rsv2 == nil {
		return errors.New("token cache is nil")
	}
	la.Cache.Set(BindUserPrefix+token, rsv2, getExpire(rsv2.loginType()))
	la.Cache.Set(TokenPrefix+token, rsv2, getExpire(rsv2.loginType()))
	return nil
}

func (la *LocalAuth) GetClaims(token string) (*Claims, error) {
	sKey := TokenPrefix + token
	if food, found := la.Cache.Get(sKey); !found || food == nil {
		return nil, fmt.Errorf("token not found:%w", ErrTokenInvalid)
	} else {
		return food.(*Claims), nil
	}
}

// Token
func (la *LocalAuth) Token(cla *Claims) (string, error) {
	userTokens, err := la.getUserTokens(cla.roleType(), cla.Id)
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
			cla.RoleType == existCla.RoleType &&
			cla.AuthId == existCla.AuthId &&
			cla.LoginType == existCla.LoginType {
			return token, nil
		}
	}
	return "", nil
}

// getUserTokens
func (la *LocalAuth) getUserTokens(roleType RoleType, userId string) (tokens, error) {
	if utokens, ok := la.Cache.Get(getPrefixKey(roleType, userId)); ok && utokens != nil {
		return utokens.(tokens), nil
	}
	return nil, nil
}

// getMultiClaimses 获取用户信息
func (la *LocalAuth) getMultiClaimses(tokens tokens) (map[string]*Claims, error) {
	clas := make(map[string]*Claims, la.getUserTokenMaxCount())
	for _, token := range tokens {
		cla, err := la.GetClaims(token)
		if err != nil {
			continue
		}
		clas[token] = cla
	}

	return clas, nil
}

func (la *LocalAuth) isUserTokenOver(roleType RoleType, userId string) bool {
	return la.getUserTokenCount(roleType, userId) >= la.getUserTokenMaxCount()
}

// getUserTokenCount 获取登录数量
func (la *LocalAuth) getUserTokenCount(roleType RoleType, userId string) int64 {
	return la.checkMaxCount(roleType, userId)
}

func (la *LocalAuth) checkMaxCount(roleType RoleType, userId string) int64 {
	utokens, _ := la.getUserTokens(roleType, userId)
	if utokens == nil {
		return 0
	}
	for index, u := range utokens {
		if _, found := la.Cache.Get(TokenPrefix + u); !found {
			if len(utokens) == 1 {
				utokens = nil
			} else {
				utokens = append(utokens[0:index], utokens[index:]...)
			}
		}
	}
	la.Cache.Set(getPrefixKey(roleType, userId), utokens, cache.NoExpiration)
	return int64(len(utokens))

}

// getUserTokenMaxCount
func (la *LocalAuth) getUserTokenMaxCount() int64 {
	if count, found := la.Cache.Get(LimitTokenPrefix); !found {
		return LimitTokenDefault
	} else {
		return count.(int64)
	}
}

// SetLimit
func (la *LocalAuth) SetLimit(tokenMaxCount int64) error {
	la.Cache.Set(LimitTokenPrefix, tokenMaxCount, cache.NoExpiration)
	return nil
}

// CleanCache
func (la *LocalAuth) CleanCache(roleType RoleType, userId string) error {
	utokens, _ := la.getUserTokens(roleType, userId)
	if utokens == nil {
		return nil
	}
	for _, token := range utokens {
		err := la.delTokenCache(token)
		if err != nil {
			continue
		}
	}
	la.Cache.Delete(getPrefixKey(roleType, userId))
	return nil
}

// IsRole
func (la *LocalAuth) IsRole(token string, roleType RoleType) (bool, error) {
	rcc, err := la.GetClaims(token)
	if err != nil {
		return false, fmt.Errorf("local: get multi claims fail %w", err)
	}
	return rcc.roleType() == roleType, nil
}

// IsRole
func (la *LocalAuth) IsSuperAdmin(token string) bool {
	rcc, err := la.GetClaims(token)
	if err != nil {
		return false
	}
	return rcc.SuperAdmin
}

func (la *LocalAuth) Close() {}
