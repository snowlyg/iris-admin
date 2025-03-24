package auth2

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

var (
	localAuth    = NewLocalAuth()
	tToken       = "TVRReU1EVTFOek13TmpFd09UWXlPRFF4TmcuTWpBeU1TMHdOeTB5T1ZRd09Ub3pNRG95T1Nzd09Eb3dNQQ.MTQyMDU1NzMwNjEwOTYyODQxNg"
	customClaims = New(
		&Multi{
			Id:            uint(1),
			SuperAdmin:    true,
			Username:      "username",
			AuthorityIds:  []string{"999"},
			AuthorityType: AdminAuthority,
			LoginType:     LoginTypeWeb,
			AuthType:      LoginTypeWeb,
			ExpiresAt:     time.Now().Local().Add(RedisSessionTimeoutWeb).Unix(),
		},
	)
	userKey = GetUserPrefixKey(customClaims.AuthorityType, customClaims.Id)
)

func TestNewLocalAuth(t *testing.T) {
	t.Run("test new local auth", func(t *testing.T) {
		if NewLocalAuth() == nil {
			t.Error("new local auth get nil")
		}
	})
}

func TestGenerateToken(t *testing.T) {
	t.Run("test generate token", func(t *testing.T) {
		token, expiresIn, err := localAuth.GenerateToken(customClaims)
		if err != nil {
			t.Fatalf("generate token %v", err)
		}
		if token == "" {
			t.Error("generate token is empty")
		}

		if expiresIn != customClaims.ExpiresAt {
			t.Errorf("generate token expires want %v but get %v", customClaims.ExpiresAt, expiresIn)
		}
		cc, err := localAuth.GetMultiClaims(token)
		if err != nil {
			t.Fatalf("get custom claims  %v", err)
		}

		if cc.Id != customClaims.Id {
			t.Errorf("get custom id want %v but get %v", customClaims.Id, cc.Id)
		}
		if cc.Username != customClaims.Username {
			t.Errorf("get custom username want %v but get %v", customClaims.Username, cc.Username)
		}
		if cc.AuthorityId != customClaims.AuthorityId {
			t.Errorf("get custom authority_id want %v but get %v", customClaims.AuthorityId, cc.AuthorityId)
		}
		if cc.AuthorityType != customClaims.AuthorityType {
			t.Errorf("get custom authority_type want %v but get %v", customClaims.AuthorityType, cc.AuthorityType)
		}
		if cc.LoginType != customClaims.LoginType {
			t.Errorf("get custom login_type want %v but get %v", customClaims.LoginType, cc.LoginType)
		}
		if cc.AuthType != customClaims.AuthType {
			t.Errorf("get custom auth_type want %v but get %v", customClaims.AuthType, cc.AuthType)
		}
		if cc.CreationTime != customClaims.CreationTime {
			t.Errorf("get custom creation_data want %v but get %v", customClaims.CreationTime, cc.CreationTime)
		}
		if cc.ExpiresAt != customClaims.ExpiresAt {
			t.Errorf("get custom expires_at want %v but get %v", customClaims.ExpiresAt, cc.ExpiresAt)
		}

		if uTokens, uFound := localAuth.Cache.Get(userKey); uFound {
			ts := uTokens.(tokens)
			if len(ts) == 0 || ts[0] != token {
				t.Errorf("user prefix value want %v but get %v", userKey, uTokens)
			}
		} else {
			t.Error("user prefix value is emptpy")
		}
		bindKey := GtSessionBindUserPrefix + token
		if uTokens, uFound := localAuth.Cache.Get(bindKey); uFound {
			if uTokens != userKey {
				t.Errorf("bind user prefix value want %v but get %v", userKey, uTokens)
			}
		} else {
			t.Error("bind user prefix value is emptpy")
		}
	})
}

func TestToCache(t *testing.T) {
	t.Run("test to cache", func(t *testing.T) {
		err := localAuth.toCache(tToken, customClaims)
		if err != nil {
			t.Fatalf("generate token %v", err)
		}
		cc, err := localAuth.GetMultiClaims(tToken)
		if err != nil {
			t.Fatalf("get custom claims  %v", err)
		}

		if cc.Id != customClaims.Id {
			t.Errorf("get custom id want %v but get %v", customClaims.Id, cc.Id)
		}
		if cc.Username != customClaims.Username {
			t.Errorf("get custom username want %v but get %v", customClaims.Username, cc.Username)
		}
		if cc.AuthorityId != customClaims.AuthorityId {
			t.Errorf("get custom authority_id want %v but get %v", customClaims.AuthorityId, cc.AuthorityId)
		}
		if cc.AuthorityType != customClaims.AuthorityType {
			t.Errorf("get custom authority_type want %v but get %v", customClaims.AuthorityType, cc.AuthorityType)
		}
		if cc.LoginType != customClaims.LoginType {
			t.Errorf("get custom login_type want %v but get %v", customClaims.LoginType, cc.LoginType)
		}
		if cc.AuthType != customClaims.AuthType {
			t.Errorf("get custom auth_type want %v but get %v", customClaims.AuthType, cc.AuthType)
		}
		if cc.CreationTime != customClaims.CreationTime {
			t.Errorf("get custom creation_data want %v but get %v", customClaims.CreationTime, cc.CreationTime)
		}
		if cc.ExpiresAt != customClaims.ExpiresAt {
			t.Errorf("get custom expires_at want %v but get %v", customClaims.ExpiresAt, cc.ExpiresAt)
		}
	})
}

func TestDelUserTokenCache(t *testing.T) {
	cc := New(
		&Multi{
			Id:            uint(2),
			Username:      "username",
			SuperAdmin:    true,
			AuthorityIds:  []string{"999"},
			AuthorityType: AdminAuthority,
			LoginType:     LoginTypeWeb,
			AuthType:      LoginTypeWeb,
			ExpiresAt:     time.Now().Local().Add(RedisSessionTimeoutWeb).Unix(),
		},
	)
	t.Run("test del user token", func(t *testing.T) {
		token, _, _ := localAuth.GenerateToken(cc)
		if token == "" {
			t.Error("generate token is empty")
		}

		err := localAuth.DelUserTokenCache(token)
		if err != nil {
			t.Fatalf("del user token cache  %v", err)
		}
		_, err = localAuth.GetMultiClaims(token)
		if !errors.Is(err, ErrTokenInvalid) {
			t.Fatalf("get custom claims err want %v but get  %v", ErrTokenInvalid, err)
		}

		if uTokens, uFound := localAuth.Cache.Get(GtSessionUserPrefix + cc.Id); uFound && uTokens != nil {
			t.Errorf("user prefix value want empty but get %v", uTokens)
		}
		bindKey := GtSessionBindUserPrefix + token
		if key, uFound := localAuth.Cache.Get(bindKey); uFound {
			t.Errorf("bind user prefix value want empty but get %v", key)
		}
	})
}

func TestIsUserTokenOver(t *testing.T) {
	cc := New(
		&Multi{
			Id:            uint(3),
			Username:      "username",
			SuperAdmin:    true,
			AuthorityIds:  []string{"999"},
			AuthorityType: AdminAuthority,
			LoginType:     LoginTypeWeb,
			AuthType:      LoginTypeWeb,
			ExpiresAt:     time.Now().Local().Add(RedisSessionTimeoutWeb).Unix(),
		},
	)
	for i := 0; i < 6; i++ {
		localAuth.GenerateToken(cc)
	}
	t.Run("test is user token over", func(t *testing.T) {
		if localAuth.isUserTokenOver(cc.AuthorityType, cc.Id) {
			t.Error("user token want not over  but get over")
		}
		count := localAuth.getUserTokenCount(cc.AuthorityType, cc.Id)
		if count != 6 {
			t.Errorf("user token count want %v  but get %v", 6, count)
		}
	})
}

func TestSetUserTokenMaxCount(t *testing.T) {
	for i := 0; i < 6; i++ {
		localAuth.GenerateToken(customClaims)
	}
	t.Run("testset user token max count", func(t *testing.T) {
		if err := localAuth.SetUserTokenMaxCount(5); err != nil {
			t.Fatalf("set user token max count %v", err)
		}
		count := localAuth.getUserTokenMaxCount()
		if count != 5 {
			t.Errorf("user token max count want %v  but get %v", 5, count)
		}
		if !localAuth.isUserTokenOver(customClaims.AuthorityType, customClaims.Id) {
			t.Error("user token want over but get not over")
		}
	})
}
func TestCleanUserTokenCache(t *testing.T) {
	for i := 0; i < 6; i++ {
		localAuth.GenerateToken(customClaims)
	}
	t.Run("test clean user token cache", func(t *testing.T) {
		if err := localAuth.CleanUserTokenCache(customClaims.AuthorityType, customClaims.Id); err != nil {
			t.Fatalf("clear user token cache %v", err)
		}
		if localAuth.getUserTokenCount(customClaims.AuthorityType, customClaims.Id) != 0 {
			t.Error("user token count want 0 but get not 0")
		}
	})
}

func TestLocalGetMultiClaims(t *testing.T) {
	defer localAuth.CleanUserTokenCache(redisClaims.AuthorityType, redisClaims.Id)
	var token string
	redisClaims.LoginType = 3
	token, _, err := localAuth.GenerateToken(redisClaims)
	if err != nil {
		t.Fatalf("get custom claims  %v", err)
	}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		redisClaims.LoginType = i
		go func(i int) {
			localAuth.GenerateToken(redisClaims)
			wg.Done()
		}(i)
		wg.Wait()
	}
	t.Run("test get custom claims", func(t *testing.T) {
		for i := 0; i < 4; i++ {
			go func() {
				cc, err := localAuth.GetMultiClaims(token)
				if err != nil {
					t.Errorf("get custom claims  %v", err)
				}
				fmt.Printf("test check token hash get %+v\n", cc)
			}()
		}
		time.Sleep(3 * time.Second)
	})
}

func TestLocalGetUserTokens(t *testing.T) {
	cc := New(
		&Multi{
			Id:            uint(121321),
			Username:      "username",
			SuperAdmin:    true,
			AuthorityIds:  []string{"999"},
			AuthorityType: AdminAuthority,
			LoginType:     LoginTypeWeb,
			AuthType:      LoginTypeWeb,
			ExpiresAt:     time.Now().Local().Add(RedisSessionTimeoutWeb).Unix(),
		},
	)
	defer localAuth.CleanUserTokenCache(cc.AuthorityType, cc.Id)
	defer localAuth.CleanUserTokenCache(redisClaims.AuthorityType, redisClaims.Id)
	token, _, err := localAuth.GenerateToken(redisClaims)
	if err != nil {
		t.Fatalf("get user tokens by claims generate token %v \n", err)
	}

	if token == "" {
		t.Fatal("get user tokens by claims generate token is empty \n")
	}

	token3232, _, err := localAuth.GenerateToken(cc)
	if err != nil {
		t.Fatalf("get user tokens by claims generate token %v \n", err)
	}

	if token3232 == "" {
		t.Fatal("get user tokens by claims generate token is empty \n")
	}

	t.Run("test get user tokens by claims", func(t *testing.T) {
		tokens, err := localAuth.getUserTokens(redisClaims.AuthorityType, redisClaims.Id)
		if err != nil {
			t.Fatalf("get user tokens by claims %v", err)
		}

		if len(tokens) != 2 {
			t.Fatalf("get user tokens by claims want len 2 but get %d", len(tokens))
		}
	})
}

func TestLocalGetTokenByClaims(t *testing.T) {
	cc := New(
		&Multi{
			Id:            uint(3232),
			Username:      "username",
			SuperAdmin:    true,
			AuthorityIds:  []string{"999"},
			AuthorityType: AdminAuthority,
			LoginType:     LoginTypeWeb,
			AuthType:      LoginTypeWeb,
			ExpiresAt:     time.Now().Local().Add(RedisSessionTimeoutWeb).Unix(),
		},
	)
	defer localAuth.CleanUserTokenCache(cc.AuthorityType, cc.Id)
	defer localAuth.CleanUserTokenCache(redisClaims.AuthorityType, redisClaims.Id)
	token, _, err := localAuth.GenerateToken(redisClaims)
	if err != nil {
		t.Fatalf("get token by claims generate token %v \n", err)
	}

	if token == "" {
		t.Fatal("get token by claims generate token is empty \n")
	}

	token3232, _, err := localAuth.GenerateToken(cc)
	if err != nil {
		t.Fatalf("get token by claims generate token %v \n", err)
	}

	if token3232 == "" {
		t.Fatal("get token by claims generate token is empty \n")
	}

	t.Run("test get token by claims", func(t *testing.T) {
		userToken, err := localAuth.GetTokenByClaims(redisClaims)
		if err != nil {
			t.Fatalf("get token by claims %v", err)
		}

		if token != userToken {
			t.Errorf("get token by claims token want %s but get %s", token, userToken)
		}
		if token == token3232 {
			t.Errorf("get token by claims token not want %s but get %s", token3232, token)
		}
	})

}
func TestLocalGetMultiClaimses(t *testing.T) {
	defer localAuth.CleanUserTokenCache(redisClaims.AuthorityType, redisClaims.Id)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		redisClaims.LoginType = i
		go func(i int) {
			localAuth.GenerateToken(redisClaims)
			wg.Done()
		}(i)
		wg.Wait()
	}
	userTokens, err := localAuth.getUserTokens(redisClaims.AuthorityType, redisClaims.Id)
	if err != nil {
		t.Fatal("get custom claimses generate token is empty \n")
	}
	t.Run("test get custom claimses", func(t *testing.T) {
		clas, err := localAuth.getMultiClaimses(userTokens)
		if err != nil {
			t.Fatalf("get custom claimses %v", err)
		}

		if len(userTokens) != 2 {
			t.Fatalf("get custom claimses want len 2 but get %d", len(userTokens))
		}
		if len(clas) != 2 {
			t.Fatalf("get custom claimses want len 2 but get %d", len(clas))
		}
	})

}
