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
	if client == nil {
		return nil, errors.New("redis client is nil")
	}
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return &RedisAuth{
		Client: client,
	}, nil
}

// Generate
func (ra *RedisAuth) Generate(claims *Claims) (string, int64, error) {
	token, err := ra.Token(claims)
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

// toCache
func (ra *RedisAuth) toCache(token string, cla *Claims) error {
	sKey := TokenPrefix + token
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

// Token
func (ra *RedisAuth) Token(cla *Claims) (string, error) {
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

// getMultiClaimses
func (ra *RedisAuth) getMultiClaimses(tokens []string) (map[string]*Claims, error) {
	clas := make(map[string]*Claims, ra.getUserTokenLimit())
	for _, token := range tokens {
		cla, err := ra.GetClaims(token)
		if err != nil {
			continue
		}
		clas[token] = cla
	}

	return clas, nil
}

// GetClaims
func (ra *RedisAuth) GetClaims(token string) (*Claims, error) {
	cla := new(Claims)
	if err := ra.Client.HGetAll(context.Background(), TokenPrefix+token).Scan(cla); err != nil {
		return nil, fmt.Errorf("get custom claims redis hgetall %w", err)
	}

	if cla.Id == "" {
		return nil, ErrEmptyToken
	}

	return cla, nil
}

// isUserTokenOver
func (ra *RedisAuth) isUserTokenOver(roleType RoleType, userId string) (bool, error) {
	max, err := ra.getUserTokenCount(roleType, userId)
	if err != nil {
		return true, err
	}
	return max >= ra.getUserTokenLimit(), nil
}

// getUserTokens
func (ra *RedisAuth) getUserTokens(roleType RoleType, userId string) ([]string, error) {
	userTokens, err := ra.Client.SMembers(context.Background(), getPrefixKey(roleType, userId)).Result()
	if err != nil {
		return nil, fmt.Errorf("get user token count menbers  %w", err)
	}
	return userTokens, nil
}

// getUserTokenCount
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

// checkUserTokenCount
func (ra *RedisAuth) checkUserTokenCount(token, userPrefixKey string) int64 {
	mun, err := ra.Client.Exists(context.Background(), TokenPrefix+token).Result()
	if err != nil || mun == 0 {
		ra.Client.SRem(context.Background(), userPrefixKey, token)
	}
	return mun
}

// getUserTokenLimit
func (ra *RedisAuth) getUserTokenLimit() int64 {
	count, err := ra.Client.Get(context.Background(), LimitTokenPrefix).Int64()
	if err != nil {
		return LimitTokenDefault
	}
	return count
}

// SetLimit
func (ra *RedisAuth) SetLimit(limit int64) error {
	err := ra.Client.Set(context.Background(), LimitTokenPrefix, limit, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// syncUserTokenCache
func (ra *RedisAuth) syncUserTokenCache(token string) error {
	cla, err := ra.GetClaims(token)
	if err != nil {
		return fmt.Errorf("sysnc user token cache %w", err)
	}
	userPrefixKey := getPrefixKey(cla.roleType(), cla.Id)
	if _, err := ra.Client.SAdd(context.Background(), userPrefixKey, token).Result(); err != nil {
		return fmt.Errorf("sync user token cache redis sadd %w", err)
	}

	bindUserPrefixKey := BindUserPrefix + token
	_, err = ra.Client.Set(context.Background(), bindUserPrefixKey, userPrefixKey, getExpire(cla.loginType())).Result()
	if err != nil {
		return fmt.Errorf("sync user token cache %w", err)
	}
	return nil
}

// UpdateCacheExpire
func (ra *RedisAuth) UpdateCacheExpire(token string) error {
	rcc, err := ra.GetClaims(token)
	if err != nil {
		return fmt.Errorf("update user token cache expire %w", err)
	}
	if rcc == nil {
		return errors.New("token cache is nil")
	}
	if err = ra.setExpire(TokenPrefix+token, rcc.loginType()); err != nil {
		return fmt.Errorf("update user token cache expire redis expire %w", err)
	}
	if err = ra.setExpire(BindUserPrefix+token, rcc.loginType()); err != nil {
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

// DelCache
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

	if e := ra.delTokenCache(token); e != nil {
		return e
	}
	return nil
}

// delUserTokenPrefixToken
func (ra *RedisAuth) delUserTokenPrefixToken(roleType RoleType, id, token string) error {
	_, err := ra.Client.SRem(context.Background(), getPrefixKey(roleType, id), token).Result()
	if err != nil {
		return fmt.Errorf("del user token cache redis srem %w", err)
	}
	return nil
}

// delTokenCache
func (ra *RedisAuth) delTokenCache(token string) error {
	sKey2 := BindUserPrefix + token
	_, err := ra.Client.Del(context.Background(), sKey2).Result()
	if err != nil {
		return fmt.Errorf("del user token cache redis del2  %w", err)
	}

	sKey3 := TokenPrefix + token
	_, err = ra.Client.Del(context.Background(), sKey3).Result()
	if err != nil {
		return fmt.Errorf("del user token cache redis del3  %w", err)
	}

	return nil
}

// CleanCache
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
