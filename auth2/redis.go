package auth2

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// RedisAuth
type RedisAuth struct {
	Client redis.UniversalClient
}

// NewRedis
func NewRedis(client redis.UniversalClient) (*RedisAuth, error) {
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("redis client is nil")
	}
	return &RedisAuth{
		Client: client,
	}, nil
}

// Generate
func (ra *RedisAuth) Generate(claims *Claims) (string, int64, error) {

	token, err := ra.Get(claims)
	if err != nil {
		return "", int64(claims.ExpiresAt), err
	}

	if token == "" {
		if isOver, err := ra.isUserTokenOver(claims.roleType(), claims.Id); err != nil {
			return "", int64(claims.ExpiresAt), err
		} else if isOver {
			return "", int64(claims.ExpiresAt), ErrOverLimit
		}

		token, err = getToken()
		if err != nil {
			return "", int64(claims.ExpiresAt), err
		}
	}

	if err = ra.toCache(token, claims); err != nil {
		return "", int64(claims.ExpiresAt), err
	}

	if err = ra.syncUserTokenCache(token); err != nil {
		return "", int64(claims.ExpiresAt), err
	}

	return token, int64(claims.ExpiresAt), nil
}

// toCache 缓存 token
func (ra *RedisAuth) toCache(token string, cla *Claims) error {
	sKey := GtSessionTokenPrefix + token
	if _, err := ra.Client.HMSet(context.Background(), sKey,
		"id", cla.Id,
		"super_admin", cla.SuperAdmin,
		"login_type", cla.LoginType,
		"auth_type", cla.AuthType,
		"username", cla.Username,
		"auth_id", cla.AuthId,
		"role_type", cla.RoleType,
		"creation_data", cla.CreationTime,
		"expires_at", cla.ExpiresAt,
	).Result(); err != nil {
		return fmt.Errorf("to cache token %w", err)
	}
	err := ra.setExpire(sKey, cla.loginType())
	if err != nil {
		return err
	}

	return nil
}

// Get 获取用户信息
func (ra *RedisAuth) Get(cla *Claims) (string, error) {
	userTokens, err := ra.getUserTokens(cla.roleType(), cla.Id)
	if err != nil {
		return "", err
	}
	clas, err := ra.getMultiClaimses(userTokens)
	if err != nil {
		return "", err
	}
	for token, existCla := range clas {
		if cla.AuthType == existCla.AuthType && cla.Id == existCla.Id && cla.RoleType == existCla.RoleType &&
			cla.AuthId == existCla.AuthId && cla.LoginType == existCla.LoginType {
			return token, nil
		}
	}
	return "", nil
}

// getMultiClaimses 获取用户信息
func (ra *RedisAuth) getMultiClaimses(tokens []string) (map[string]*Claims, error) {
	clas := make(map[string]*Claims, ra.getUserTokenMaxCount())
	for _, token := range tokens {
		cla, err := ra.GetClaims(token)
		if err != nil {
			continue
		}
		clas[token] = cla
	}

	return clas, nil
}

// GetClaims 获取用户信息
func (ra *RedisAuth) GetClaims(token string) (*Claims, error) {
	cla := new(Claims)
	if err := ra.Client.HGetAll(context.Background(), GtSessionTokenPrefix+token).Scan(cla); err != nil {
		return nil, fmt.Errorf("get custom claims redis hgetall %w", err)
	}

	if cla.Id == "" {
		return nil, ErrEmptyToken
	}

	return cla, nil
}

// isUserTokenOver 超过登录设备限制
func (ra *RedisAuth) isUserTokenOver(roleType RoleType, userId string) (bool, error) {
	max, err := ra.getUserTokenCount(roleType, userId)
	if err != nil {
		return true, err
	}
	return max >= ra.getUserTokenMaxCount(), nil
}

// getUserTokens 获取登录数量
func (ra *RedisAuth) getUserTokens(roleType RoleType, userId string) ([]string, error) {
	userTokens, err := ra.Client.SMembers(context.Background(), getPrefixKey(roleType, userId)).Result()
	if err != nil {
		return nil, fmt.Errorf("get user token count menbers  %w", err)
	}
	return userTokens, nil
}

// getUserTokenCount 获取登录数量
func (ra *RedisAuth) getUserTokenCount(roleType RoleType, userId string) (int64, error) {
	var count int64
	userTokens, err := ra.getUserTokens(roleType, userId)
	if err != nil {
		return count, fmt.Errorf("get user token count menbers  %w", err)
	}
	userPrefixKey := getPrefixKey(roleType, userId)
	for _, token := range userTokens {
		if ra.checkUserTokenCount(token, userPrefixKey) == 1 {
			count++
		}
	}
	return count, nil
}

// checkUserTokenCount 验证登录数量,清除 userPrefixKey 下无效 token
func (ra *RedisAuth) checkUserTokenCount(token, userPrefixKey string) int64 {
	mun, err := ra.Client.Exists(context.Background(), GtSessionTokenPrefix+token).Result()
	if err != nil || mun == 0 {
		ra.Client.SRem(context.Background(), userPrefixKey, token)
	}
	return mun
}

// getUserTokenMaxCount 最大登录限制
func (ra *RedisAuth) getUserTokenMaxCount() int64 {
	count, err := ra.Client.Get(context.Background(), GtSessionUserMaxTokenPrefix).Int64()
	if err != nil {
		return GtSessionUserMaxTokenDefault
	}
	return count
}

// SetMaxCount 最大登录限制
func (ra *RedisAuth) SetMaxCount(tokenMaxCount int64) error {
	err := ra.Client.Set(context.Background(), GtSessionUserMaxTokenPrefix, tokenMaxCount, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// syncUserTokenCache 同步 token 到用户缓存
func (ra *RedisAuth) syncUserTokenCache(token string) error {
	cla, err := ra.GetClaims(token)
	if err != nil {
		return fmt.Errorf("sysnc user token cache %w", err)
	}
	userPrefixKey := getPrefixKey(cla.roleType(), cla.Id)
	if _, err := ra.Client.SAdd(context.Background(), userPrefixKey, token).Result(); err != nil {
		return fmt.Errorf("sync user token cache redis sadd %w", err)
	}

	bindUserPrefixKey := GtSessionBindUserPrefix + token
	_, err = ra.Client.Set(context.Background(), bindUserPrefixKey, userPrefixKey, getExpire(cla.loginType())).Result()
	if err != nil {
		return fmt.Errorf("sync user token cache %w", err)
	}
	return nil
}

// UpdateCacheExpire 更新过期时间
func (ra *RedisAuth) UpdateCacheExpire(token string) error {
	rcc, err := ra.GetClaims(token)
	if err != nil {
		return fmt.Errorf("update user token cache expire %w", err)
	}
	if rcc == nil {
		return errors.New("token cache is nil")
	}
	if err = ra.setExpire(GtSessionTokenPrefix+token, rcc.loginType()); err != nil {
		return fmt.Errorf("update user token cache expire redis expire %w", err)
	}
	if err = ra.setExpire(GtSessionBindUserPrefix+token, rcc.loginType()); err != nil {
		return fmt.Errorf("update user token cache expire redis expire %w", err)
	}
	return nil
}

func (ra *RedisAuth) setExpire(key string, loginType LoginType) error {
	if _, err := ra.Client.Expire(context.Background(), key, getExpire(loginType)).Result(); err != nil {
		return fmt.Errorf("update user token cache expire redis expire %w", err)
	}
	return nil
}

// DelCache 删除token缓存
func (ra *RedisAuth) DelCache(token string) error {
	log.Println("auth2: redis del user token")
	cla, err := ra.GetClaims(token)
	if err != nil {
		return err
	}
	if cla == nil {
		return errors.New("del user token, reids cache is nil")
	}

	if e := ra.delUserTokenPrefixToken(cla.roleType(), cla.Id, token); e != nil {
		return e
	}

	err = ra.delTokenCache(token)
	if err != nil {
		return err
	}
	return nil
}

// delUserTokenPrefixToken 删除 user token缓存
func (ra *RedisAuth) delUserTokenPrefixToken(roleType RoleType, id, token string) error {
	_, err := ra.Client.SRem(context.Background(), getPrefixKey(roleType, id), token).Result()
	if err != nil {
		return fmt.Errorf("del user token cache redis srem %w", err)
	}
	return nil
}

// delTokenCache 删除token缓存
func (ra *RedisAuth) delTokenCache(token string) error {
	sKey2 := GtSessionBindUserPrefix + token
	_, err := ra.Client.Del(context.Background(), sKey2).Result()
	if err != nil {
		return fmt.Errorf("del user token cache redis del2  %w", err)
	}

	sKey3 := GtSessionTokenPrefix + token
	_, err = ra.Client.Del(context.Background(), sKey3).Result()
	if err != nil {
		return fmt.Errorf("del user token cache redis del3  %w", err)
	}

	return nil
}

// CleanCache 清空token缓存
func (ra *RedisAuth) CleanCache(roleType RoleType, userId string) error {
	allTokens, err := ra.getUserTokens(roleType, userId)
	if err != nil {
		return fmt.Errorf("clean user token cache redis smembers  %w", err)
	}
	_, err = ra.Client.Del(context.Background(), getPrefixKey(roleType, userId)).Result()
	if err != nil {
		return fmt.Errorf("clean user token cache redis del  %w", err)
	}

	for _, token := range allTokens {
		err = ra.delTokenCache(token)
		if err != nil {
			return err
		}
	}
	return nil
}

// IsRole
func (ra *RedisAuth) IsRole(token string, roleType RoleType) (bool, error) {
	rcc, err := ra.GetClaims(token)
	if err != nil {
		return false, fmt.Errorf("get User's infomation return error: %w", err)
	}
	return rcc.roleType() == roleType, nil
}

// IsSuperAdmin
func (ra *RedisAuth) IsSuperAdmin(token string) bool {
	rcc, err := ra.GetClaims(token)
	if err != nil {
		return false
	}
	return rcc.SuperAdmin
}

// Close
func (ra *RedisAuth) Close() {
	ra.Client.Close()
}
