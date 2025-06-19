package auth2

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	wg      sync.WaitGroup
	options = &redis.UniversalOptions{
		DB:          1,
		Addrs:       []string{os.Getenv("IRIS_ADMIN_REDIS_ADDR")},
		Password:    os.Getenv("IRIS_ADMIN_REDIS_PWD"),
		PoolSize:    10,
		IdleTimeout: 300 * time.Second,
		// Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
		// 	conn, err := net.Dial(network, addr)
		// 	if err == nil {
		// 		go func() {
		// 			time.Sleep(5 * time.Second)
		// 			conn.Close()
		// 		}()
		// 	}
		// 	return conn, err
		// },
	}

	rToken     = "TVRReU1EVTFOek13TmpFd09UWXlPRFF4TmcuTWpBeU1TMHdOeTB5T1ZRd09Ub3pNRG95T1Nzd09Eb3dNQQ.MTQyMDU1NzMwNjEwOTYyODrtrt"
	logTypeWeb = NewClaims(
		&Agent{
			Id:         uint(121321),
			Username:   "username",
			SuperAdmin: true,
			AuthIds:    []string{"999"},
			RoleType:   RoleAdmin,
			LoginType:  LoginTypeWeb,
			AuthType:   AuthPwd,
			ExpiresAt:  time.Now().Local().Add(TimeoutWeb).Unix(),
		},
	)
	ruserKey = getPrefixKey(logTypeWeb.roleType(), logTypeWeb.Id)
)

func TestRedisGenerateToken(t *testing.T) {
	redisPwd := os.Getenv("IRIS_ADMIN_REDIS_PWD")
	if redisPwd == "" {
		t.SkipNow()
	}
	redisAuth, err := NewRedis(redis.NewUniversalClient(options))
	if err != nil {
		t.Fatal(err.Error())
	}
	defer redisAuth.CleanCache(logTypeWeb.roleType(), logTypeWeb.Id)
	token, expiresIn, err := redisAuth.Generate(logTypeWeb)
	if err != nil {
		t.Fatalf("generate token %v", err)
	}
	if token == "" {
		t.Error("generate token is empty")
	}

	if expiresIn != logTypeWeb.ExpiresAt {
		t.Errorf("generate token expires want %v but get %v", logTypeWeb.ExpiresAt, expiresIn)
	}
	cc, err := redisAuth.GetClaims(token)
	if err != nil {
		t.Fatalf("get custom claims  %v", err)
	}

	if cc.Id != logTypeWeb.Id {
		t.Errorf("get custom id want %v but get %v", logTypeWeb.Id, cc.Id)
	}
	if cc.Username != logTypeWeb.Username {
		t.Errorf("get custom username want %v but get %v", logTypeWeb.Username, cc.Username)
	}
	if cc.AuthId != logTypeWeb.AuthId {
		t.Errorf("get custom authority_id want %v but get %v", logTypeWeb.AuthId, cc.AuthId)
	}
	if cc.RoleType != logTypeWeb.RoleType {
		t.Errorf("get custom authority_type want %v but get %v", logTypeWeb.RoleType, cc.RoleType)
	}
	if cc.LoginType != logTypeWeb.LoginType {
		t.Errorf("get custom login_type want %v but get %v", logTypeWeb.LoginType, cc.LoginType)
	}
	if cc.AuthType != logTypeWeb.AuthType {
		t.Errorf("get custom auth_type want %v but get %v", logTypeWeb.AuthType, cc.AuthType)
	}
	if cc.CreationTime != logTypeWeb.CreationTime {
		t.Errorf("get custom creation_data want %v but get %v", logTypeWeb.CreationTime, cc.CreationTime)
	}
	if cc.ExpiresAt != logTypeWeb.ExpiresAt {
		t.Errorf("get custom expires_at want %v but get %v", logTypeWeb.ExpiresAt, cc.ExpiresAt)
	}

	if uTokens, err := redisAuth.Client.SMembers(context.Background(), ruserKey).Result(); err != nil {
		t.Fatalf("user prefix value get %s", err)
	} else {
		if len(uTokens) == 0 || uTokens[0] != token {
			t.Errorf("user prefix value want %v but get %v", ruserKey, uTokens)
		}
	}
	bindKey := BindUserPrefix + token
	key, err := redisAuth.Client.Get(context.Background(), bindKey).Result()
	if err != nil {
		t.Fatal(err)
	}
	if key != ruserKey {
		t.Errorf("bind user prefix value want %v but get %v", ruserKey, key)
	}
}

func TestRedisToCache(t *testing.T) {
	redisPwd := os.Getenv("IRIS_ADMIN_REDIS_PWD")
	if redisPwd == "" {
		t.SkipNow()
	}
	redisAuth, err := NewRedis(redis.NewUniversalClient(options))
	if err != nil {
		t.Fatal(err.Error())
	}
	defer redisAuth.Client.Del(context.Background(), TokenPrefix+rToken)
	if err := redisAuth.toCache(rToken, logTypeWeb); err != nil {
		t.Fatalf("generate token %v", err)
	}
	cc, err := redisAuth.GetClaims(rToken)
	if err != nil {
		t.Fatalf("get custom claims  %v", err)
	}

	if cc.Id != logTypeWeb.Id {
		t.Errorf("get custom id want %v but get %v", logTypeWeb.Id, cc.Id)
	}
	if cc.Username != logTypeWeb.Username {
		t.Errorf("get custom username want %v but get %v", logTypeWeb.Username, cc.Username)
	}
	if cc.AuthId != logTypeWeb.AuthId {
		t.Errorf("get custom authority_id want %v but get %v", logTypeWeb.AuthId, cc.AuthId)
	}
	if cc.RoleType != logTypeWeb.RoleType {
		t.Errorf("get custom authority_type want %v but get %v", logTypeWeb.RoleType, cc.RoleType)
	}
	if cc.LoginType != logTypeWeb.LoginType {
		t.Errorf("get custom login_type want %v but get %v", logTypeWeb.LoginType, cc.LoginType)
	}
	if cc.AuthType != logTypeWeb.AuthType {
		t.Errorf("get custom auth_type want %v but get %v", logTypeWeb.AuthType, cc.AuthType)
	}
	if cc.CreationTime != logTypeWeb.CreationTime {
		t.Errorf("get custom creation_data want %v but get %v", logTypeWeb.CreationTime, cc.CreationTime)
	}
	if cc.ExpiresAt != logTypeWeb.ExpiresAt {
		t.Errorf("get custom expires_at want %v but get %v", logTypeWeb.ExpiresAt, cc.ExpiresAt)
	}
}

func TestRedisDelUserTokenCache(t *testing.T) {
	redisPwd := os.Getenv("IRIS_ADMIN_REDIS_PWD")
	if redisPwd == "" {
		t.SkipNow()
	}
	cc := NewClaims(
		&Agent{
			Id:         uint(221),
			Username:   "username",
			SuperAdmin: true,
			AuthIds:    []string{"999"},
			RoleType:   RoleAdmin,
			LoginType:  LoginTypeWeb,
			AuthType:   AuthPwd,
			ExpiresAt:  time.Now().Local().Add(TimeoutWeb).Unix(),
		},
	)
	redisAuth, err := NewRedis(redis.NewUniversalClient(options))
	if err != nil {
		t.Fatal(err.Error())
	}
	defer redisAuth.CleanCache(cc.roleType(), cc.Id)
	token, _, _ := redisAuth.Generate(cc)
	if token == "" {
		t.Error("generate token is empty")
	}

	if err := redisAuth.DelCache(token); err != nil {
		t.Fatalf("del user token cache  %v", err)
	}
	_, err = redisAuth.GetClaims(token)
	if !errors.Is(err, ErrEmptyToken) {
		t.Fatalf("get custom claims err want '%v' but get  '%v'", ErrEmptyToken, err)
	}

	if uTokens, err := redisAuth.Client.SMembers(context.Background(), UserPrefix+cc.Id).Result(); err != nil {
		t.Fatalf("user prefix value wantget %v", err)
	} else if len(uTokens) != 0 {
		t.Errorf("user prefix value want empty but get %+v", uTokens)
	}
	bindKey := BindUserPrefix + token
	key, _ := redisAuth.Client.Get(context.Background(), bindKey).Result()
	if key != "" {
		t.Errorf("bind user prefix value want empty but get %v", key)
	}
}

func TestRedisIsUserTokenOver(t *testing.T) {
	redisPwd := os.Getenv("IRIS_ADMIN_REDIS_PWD")
	if redisPwd == "" {
		t.SkipNow()
	}
	cc := NewClaims(
		&Agent{
			Id:         uint(3232),
			Username:   "username",
			SuperAdmin: true,
			AuthIds:    []string{"999"},
			RoleType:   RoleAdmin,
			LoginType:  LoginTypeWeb,
			AuthType:   AuthPwd,
			ExpiresAt:  time.Now().Local().Add(TimeoutWeb).Unix(),
		},
	)
	redisAuth, err := NewRedis(redis.NewUniversalClient(options))
	if err != nil {
		t.Fatal(err.Error())
	}
	defer redisAuth.CleanCache(cc.roleType(), cc.Id)
	if err := redisAuth.SetLimit(10); err != nil {
		t.Fatalf("set user token max count %v", err)
	}
	var wantTokenLen int64 = 0
	for i := LoginTypeWeb; i <= LoginTypeWx; i++ {
		cc.setLoginType(int(i))
		wg.Add(1)
		wantTokenLen++
		go func(i LoginType) {
			redisAuth.Generate(cc)
			wg.Done()
		}(i)
		wg.Wait()
	}
	isOver, err := redisAuth.isUserTokenOver(cc.roleType(), cc.Id)
	if err != nil {
		t.Fatalf("is user token over get %v", err)
	}
	if isOver {
		t.Error("user token want not over  but get over")
	}
	count, err := redisAuth.getUserTokenCount(cc.roleType(), cc.Id)
	if err != nil {
		t.Fatalf("user token count get %v", err)
	}
	if count != wantTokenLen {
		t.Errorf("user token count want %v but get %v", wantTokenLen, count)
	}
}

func TestRedisSetUserTokenMaxCount(t *testing.T) {
	redisPwd := os.Getenv("IRIS_ADMIN_REDIS_PWD")
	if redisPwd == "" {
		t.SkipNow()
	}
	redisAuth, err := NewRedis(redis.NewUniversalClient(options))
	if err != nil {
		t.Fatal(err.Error())
	}
	defer redisAuth.CleanCache(logTypeWeb.roleType(), logTypeWeb.Id)
	if err := redisAuth.SetLimit(10); err != nil {
		t.Fatalf("set user token max count %v", err)
	}
	for i := LoginTypeWeb; i <= LoginTypeWx; i++ {
		wg.Add(1)
		logTypeWeb.setLoginType(int(i))
		go func(i LoginType) {
			redisAuth.Generate(logTypeWeb)
			wg.Done()
		}(i)
		wg.Wait()
	}
	if err := redisAuth.SetLimit(3); err != nil {
		t.Fatalf("set user token max count %v", err)
	}
	count := redisAuth.getUserTokenLimit()
	if count != 3 {
		t.Errorf("user token max count want %v  but get %v", 3, count)
	}
	isOver, err := redisAuth.isUserTokenOver(logTypeWeb.roleType(), logTypeWeb.Id)
	if err != nil {
		t.Fatalf("is user token over get %v", err)
	}
	if !isOver {
		t.Error("user token want over but get not over")
	}
}
func TestRedisCleanUserTokenCache(t *testing.T) {
	redisPwd := os.Getenv("IRIS_ADMIN_REDIS_PWD")
	if redisPwd == "" {
		t.SkipNow()
	}
	redisAuth, err := NewRedis(redis.NewUniversalClient(options))
	if err != nil {
		t.Fatal(err.Error())
	}
	defer redisAuth.CleanCache(logTypeWeb.roleType(), logTypeWeb.Id)
	for i := LoginTypeWeb; i <= LoginTypeWx; i++ {
		wg.Add(1)
		logTypeWeb.setLoginType(int(i))
		go func(i LoginType) {
			redisAuth.Generate(logTypeWeb)
			wg.Done()
		}(i)
		wg.Wait()
	}
	if err := redisAuth.CleanCache(logTypeWeb.roleType(), logTypeWeb.Id); err != nil {
		t.Fatalf("clear user token cache %v", err)
	}
	count, err := redisAuth.getUserTokenCount(logTypeWeb.roleType(), logTypeWeb.Id)
	if err != nil {
		t.Fatalf("user token count get %v", err)
	}
	if count != 0 {
		t.Error("user token count want 0 but get not 0")
	}
}

func TestRedisGetMultiClaims(t *testing.T) {
	redisPwd := os.Getenv("IRIS_ADMIN_REDIS_PWD")
	if redisPwd == "" {
		t.SkipNow()
	}
	redisAuth, err := NewRedis(redis.NewUniversalClient(options))
	if err != nil {
		t.Fatal(err.Error())
	}
	defer redisAuth.CleanCache(logTypeWeb.roleType(), logTypeWeb.Id)
	logTypeWeb.LoginType = 3
	token, _, err := redisAuth.Generate(logTypeWeb)
	if err != nil {
		t.Fatalf("get custom claims  %v", err)
	}
	for i := LoginTypeWeb; i <= LoginTypeWx; i++ {
		wg.Add(1)
		logTypeWeb.setLoginType(int(i))
		go func(i LoginType) {
			redisAuth.Generate(logTypeWeb)
			wg.Done()
		}(i)
		wg.Wait()
	}
	for i := 0; i < 4; i++ {
		go func() {
			_, err := redisAuth.GetClaims(token)
			if err != nil {
				t.Errorf("get custom claims  %v", err)
			}
		}()
	}
	time.Sleep(3 * time.Second)
}

func TestRedisGetUserTokens(t *testing.T) {
	redisPwd := os.Getenv("IRIS_ADMIN_REDIS_PWD")
	if redisPwd == "" {
		t.SkipNow()
	}
	cc := NewClaims(
		&Agent{
			Id:         uint(121321),
			Username:   "username",
			SuperAdmin: true,
			AuthIds:    []string{"999"},
			RoleType:   RoleAdmin,
			LoginType:  LoginTypeWeb,
			AuthType:   AuthPwd,
			ExpiresAt:  time.Now().Local().Add(TimeoutWeb).Unix(),
		},
	)
	redisAuth, err := NewRedis(redis.NewUniversalClient(options))
	if err != nil {
		t.Fatal(err.Error())
	}
	defer redisAuth.CleanCache(cc.roleType(), cc.Id)
	defer redisAuth.CleanCache(logTypeWeb.roleType(), logTypeWeb.Id)
	token, _, err := redisAuth.Generate(logTypeWeb)
	if err != nil {
		t.Fatalf("get user tokens by claims generate token %v \n", err)
	}

	if token == "" {
		t.Fatal("get user tokens by claims generate token is empty \n")
	}

	token3232, _, err := redisAuth.Generate(cc)
	if err != nil {
		t.Fatalf("get user tokens by claims generate token %v \n", err)
	}

	if token3232 == "" {
		t.Fatal("get user tokens by claims generate token is empty \n")
	}

	tokens, err := redisAuth.getUserTokens(logTypeWeb.roleType(), logTypeWeb.Id)
	if err != nil {
		t.Fatalf("get user tokens by claims %v", err)
	}
	wantTokenLen := 2
	if len(tokens) != wantTokenLen {
		t.Fatalf("get user tokens by claims want len %d but get %d", wantTokenLen, len(tokens))
	}
}

func TestRedisGetTokenByClaims(t *testing.T) {
	redisPwd := os.Getenv("IRIS_ADMIN_REDIS_PWD")
	if redisPwd == "" {
		t.SkipNow()
	}
	cc := NewClaims(
		&Agent{
			Id:         uint(3232),
			Username:   "username",
			SuperAdmin: true,
			AuthIds:    []string{"999"},
			RoleType:   RoleAdmin,
			LoginType:  LoginTypeWeb,
			AuthType:   AuthPwd,
			ExpiresAt:  time.Now().Local().Add(TimeoutWeb).Unix(),
		},
	)
	redisAuth, err := NewRedis(redis.NewUniversalClient(options))
	if err != nil {
		t.Fatal(err.Error())
	}
	defer redisAuth.CleanCache(cc.roleType(), cc.Id)
	defer redisAuth.CleanCache(logTypeWeb.roleType(), logTypeWeb.Id)

	token, _, err := redisAuth.Generate(logTypeWeb)
	if err != nil {
		t.Fatalf("get token by claims generate token %v \n", err)
	}

	if token == "" {
		t.Fatal("get token by claims generate token is empty \n")
	}

	token3232, _, err := redisAuth.Generate(cc)
	if err != nil {
		t.Fatalf("get token by claims generate token %v \n", err)
	}

	if token3232 == "" {
		t.Fatal("get token by claims generate token is empty \n")
	}

	userToken, err := redisAuth.Token(logTypeWeb)
	if err != nil {
		t.Fatalf("get token by claims %v", err)
	}

	if token != userToken {
		t.Errorf("get token by claims token want %s but get %s", token, userToken)
	}
	if token == token3232 {
		t.Errorf("get token by claims token not want %s but get %s", token3232, token)
	}

}
func TestRedisGetMultiClaimses(t *testing.T) {
	redisPwd := os.Getenv("IRIS_ADMIN_REDIS_PWD")
	if redisPwd == "" {
		t.SkipNow()
	}
	redisAuth, err := NewRedis(redis.NewUniversalClient(options))
	if err != nil {
		t.Fatal(err.Error())
	}
	defer redisAuth.CleanCache(logTypeWeb.roleType(), logTypeWeb.Id)
	wantTokenLen := 0
	for i := LoginTypeWeb; i <= LoginTypeWx; i++ {
		wg.Add(1)
		wantTokenLen++
		logTypeWeb.setLoginType(int(i))
		go func(i LoginType) {
			redisAuth.Generate(logTypeWeb)
			wg.Done()
		}(i)
		wg.Wait()
	}
	userTokens, err := redisAuth.getUserTokens(logTypeWeb.roleType(), logTypeWeb.Id)
	if err != nil {
		t.Fatal("get custom claimses generate token is empty \n")
	}
	clas, err := redisAuth.getMultiClaimses(userTokens)
	if err != nil {
		t.Fatalf("get custom claimses %v", err)
	}

	if len(userTokens) != wantTokenLen {
		t.Fatalf("get custom claimses want len %d but get %d", wantTokenLen, len(userTokens))
	}

	if len(clas) != wantTokenLen {
		t.Fatalf("get custom claimses want len %d but get %d", wantTokenLen, len(clas))
	}

}
