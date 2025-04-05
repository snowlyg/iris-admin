package auth2

import (
	"errors"
	"testing"
	"time"
)

var (
	localAuth    = NewLocal()
	tToken       = "TVRReU1EVTFOek13TmpFd09UWXlPRFF4TmcuTWpBeU1TMHdOeTB5T1ZRd09Ub3pNRG95T1Nzd09Eb3dNQQ.MTQyMDU1NzMwNjEwOTYyODQxNg"
	loginTypeApp = NewClaims(
		&Agent{
			Id:         uint(1),
			SuperAdmin: true,
			Username:   "username",
			AuthIds:    []string{"999"},
			RoleType:   RoleAdmin,
			LoginType:  LoginTypeWeb,
			AuthType:   AuthPwd,
			ExpiresAt:  time.Now().Local().Add(RedisSessionTimeoutWeb).Unix(),
		},
	)
	userKey = getPrefixKey(loginTypeApp.roleType(), loginTypeApp.Id)
)

func TestNewLocalAuth(t *testing.T) {
	if NewLocal() == nil {
		t.Error("new local auth get nil")
	}
}

func TestGenerateToken(t *testing.T) {
	token, expiresIn, err := localAuth.Generate(loginTypeApp)
	if err != nil {
		t.Fatalf("generate token %v", err)
	}
	if token == "" {
		t.Error("generate token is empty")
	}

	if expiresIn != loginTypeApp.ExpiresAt {
		t.Errorf("generate token expires want %v but get %v", loginTypeApp.ExpiresAt, expiresIn)
	}
	cc, err := localAuth.GetClaims(token)
	if err != nil {
		t.Fatalf("get custom claims  %v", err)
	}

	if cc.Id != loginTypeApp.Id {
		t.Errorf("get custom id want %v but get %v", loginTypeApp.Id, cc.Id)
	}
	if cc.Username != loginTypeApp.Username {
		t.Errorf("get custom username want %v but get %v", loginTypeApp.Username, cc.Username)
	}
	if cc.AuthId != loginTypeApp.AuthId {
		t.Errorf("get custom authority_id want %v but get %v", loginTypeApp.AuthId, cc.AuthId)
	}
	if cc.RoleType != loginTypeApp.RoleType {
		t.Errorf("get custom authority_type want %v but get %v", loginTypeApp.RoleType, cc.RoleType)
	}
	if cc.LoginType != loginTypeApp.LoginType {
		t.Errorf("get custom login_type want %v but get %v", loginTypeApp.LoginType, cc.LoginType)
	}
	if cc.AuthType != loginTypeApp.AuthType {
		t.Errorf("get custom auth_type want %v but get %v", loginTypeApp.AuthType, cc.AuthType)
	}
	if cc.CreationTime != loginTypeApp.CreationTime {
		t.Errorf("get custom creation_data want %v but get %v", loginTypeApp.CreationTime, cc.CreationTime)
	}
	if cc.ExpiresAt != loginTypeApp.ExpiresAt {
		t.Errorf("get custom expires_at want %v but get %v", loginTypeApp.ExpiresAt, cc.ExpiresAt)
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
}

func TestToCache(t *testing.T) {
	err := localAuth.toCache(tToken, loginTypeApp)
	if err != nil {
		t.Fatalf("generate token %v", err)
	}
	cc, err := localAuth.GetClaims(tToken)
	if err != nil {
		t.Fatalf("get custom claims  %v", err)
	}

	if cc.Id != loginTypeApp.Id {
		t.Errorf("get custom id want %v but get %v", loginTypeApp.Id, cc.Id)
	}
	if cc.Username != loginTypeApp.Username {
		t.Errorf("get custom username want %v but get %v", loginTypeApp.Username, cc.Username)
	}
	if cc.AuthId != loginTypeApp.AuthId {
		t.Errorf("get custom authority_id want %v but get %v", loginTypeApp.AuthId, cc.AuthId)
	}
	if cc.RoleType != loginTypeApp.RoleType {
		t.Errorf("get custom authority_type want %v but get %v", loginTypeApp.RoleType, cc.RoleType)
	}
	if cc.LoginType != loginTypeApp.LoginType {
		t.Errorf("get custom login_type want %v but get %v", loginTypeApp.LoginType, cc.LoginType)
	}
	if cc.AuthType != loginTypeApp.AuthType {
		t.Errorf("get custom auth_type want %v but get %v", loginTypeApp.AuthType, cc.AuthType)
	}
	if cc.CreationTime != loginTypeApp.CreationTime {
		t.Errorf("get custom creation_data want %v but get %v", loginTypeApp.CreationTime, cc.CreationTime)
	}
	if cc.ExpiresAt != loginTypeApp.ExpiresAt {
		t.Errorf("get custom expires_at want %v but get %v", loginTypeApp.ExpiresAt, cc.ExpiresAt)
	}
}

func TestDelUserTokenCache(t *testing.T) {
	cc := NewClaims(
		&Agent{
			Id:         uint(2),
			Username:   "username",
			SuperAdmin: true,
			AuthIds:    []string{"999"},
			RoleType:   RoleAdmin,
			LoginType:  LoginTypeWeb,
			AuthType:   AuthPwd,
			ExpiresAt:  time.Now().Local().Add(RedisSessionTimeoutWeb).Unix(),
		},
	)
	token, _, _ := localAuth.Generate(cc)
	if token == "" {
		t.Error("generate token is empty")
	}

	err := localAuth.DelCache(token)
	if err != nil {
		t.Fatalf("del user token cache  %v", err)
	}
	_, err = localAuth.GetClaims(token)
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
}

func TestIsUserTokenOver(t *testing.T) {
	cc := NewClaims(
		&Agent{
			Id:         uint(3),
			Username:   "username",
			SuperAdmin: true,
			AuthIds:    []string{"999"},
			RoleType:   RoleAdmin,
			LoginType:  LoginTypeWeb,
			AuthType:   AuthPwd,
			ExpiresAt:  time.Now().Local().Add(RedisSessionTimeoutWeb).Unix(),
		},
	)
	for i := 0; i < 6; i++ {
		localAuth.Generate(cc)
	}
	if localAuth.isUserTokenOver(cc.roleType(), cc.Id) {
		t.Error("user token want not over  but get over")
	}
	count := localAuth.getUserTokenCount(cc.roleType(), cc.Id)
	if count != 6 {
		t.Errorf("user token count want %v  but get %v", 6, count)
	}
}

func TestSetUserTokenMaxCount(t *testing.T) {
	for i := 0; i < 6; i++ {
		localAuth.Generate(loginTypeApp)
	}
	if err := localAuth.SetMaxCount(5); err != nil {
		t.Fatalf("set user token max count %v", err)
	}
	count := localAuth.getUserTokenMaxCount()
	if count != 5 {
		t.Errorf("user token max count want %v  but get %v", 5, count)
	}
	if !localAuth.isUserTokenOver(loginTypeApp.roleType(), loginTypeApp.Id) {
		t.Error("user token want over but get not over")
	}
}
func TestCleanUserTokenCache(t *testing.T) {
	for i := 0; i < 6; i++ {
		localAuth.Generate(loginTypeApp)
	}
	if err := localAuth.CleanCache(loginTypeApp.roleType(), loginTypeApp.Id); err != nil {
		t.Fatalf("clear user token cache %v", err)
	}
	if localAuth.getUserTokenCount(loginTypeApp.roleType(), loginTypeApp.Id) != 0 {
		t.Error("user token count want 0 but get not 0")
	}
}

func TestLocalGetMultiClaims(t *testing.T) {
	defer localAuth.CleanCache(loginTypeApp.roleType(), loginTypeApp.Id)
	var token string
	loginTypeApp.LoginType = 3
	token, _, err := localAuth.Generate(loginTypeApp)
	if err != nil {
		t.Fatalf("get custom claims  %v", err)
	}
	for i := LoginTypeWeb; i <= LoginTypeDevice; i++ {
		wg.Add(1)
		loginTypeApp.setLoginType(int(i))
		go func(i LoginType) {
			localAuth.Generate(loginTypeApp)
			wg.Done()
		}(i)
		wg.Wait()
	}
	for i := 0; i < 4; i++ {
		go func() {
			_, err := localAuth.GetClaims(token)
			if err != nil {
				t.Errorf("get custom claims fail:%v", err)
			}
		}()
	}
	time.Sleep(3 * time.Second)
}

func TestLocalGetUserTokens(t *testing.T) {
	loginTypeWeb := NewClaims(
		&Agent{
			Id:         uint(121321),
			Username:   "username",
			SuperAdmin: true,
			AuthIds:    []string{"999"},
			RoleType:   RoleAdmin,
			LoginType:  LoginTypeWeb,
			AuthType:   AuthPwd,
			ExpiresAt:  time.Now().Local().Add(RedisSessionTimeoutWeb).Unix(),
		},
	)

	defer localAuth.CleanCache(loginTypeWeb.roleType(), loginTypeWeb.Id)
	defer localAuth.CleanCache(loginTypeApp.roleType(), loginTypeApp.Id)

	token, _, err := localAuth.Generate(loginTypeApp)
	if err != nil {
		t.Fatalf("get user tokens by claims generate token %v \n", err)
	}

	if token == "" {
		t.Fatal("get user tokens by claims generate token is empty \n")
	}

	token3232, _, err := localAuth.Generate(loginTypeWeb)
	if err != nil {
		t.Fatalf("get user tokens by claims generate token %v \n", err)
	}

	if token3232 == "" {
		t.Fatal("get user tokens by claims generate token is empty \n")
	}

	if token == token3232 {
		t.Fatal("get user tokens by claims generate token is same")
	}

	tokens, err := localAuth.getUserTokens(loginTypeApp.roleType(), loginTypeApp.Id)
	if err != nil {
		t.Fatalf("get user tokens by claims %v", err)
	}

	if len(tokens) != 1 {
		t.Fatalf("get user tokens by claims want len 1 but get %d", len(tokens))
	}
}

func TestLocalGetTokenByClaims(t *testing.T) {
	LoginTypeWeb := NewClaims(
		&Agent{
			Id:         uint(3232),
			Username:   "username",
			SuperAdmin: true,
			AuthIds:    []string{"999"},
			RoleType:   RoleAdmin,
			LoginType:  LoginTypeWeb,
			AuthType:   AuthPwd,
			ExpiresAt:  time.Now().Local().Add(RedisSessionTimeoutWeb).Unix(),
		},
	)
	defer localAuth.CleanCache(LoginTypeWeb.roleType(), LoginTypeWeb.Id)
	defer localAuth.CleanCache(loginTypeApp.roleType(), loginTypeApp.Id)

	token, _, err := localAuth.Generate(loginTypeApp)
	if err != nil {
		t.Fatalf("get token by claims generate token %v \n", err)
	}

	if token == "" {
		t.Fatal("get token by claims generate token is empty \n")
	}

	token3232, _, err := localAuth.Generate(LoginTypeWeb)
	if err != nil {
		t.Fatalf("get token by claims generate token %v \n", err)
	}

	if token3232 == "" {
		t.Fatal("get token by claims generate token is empty \n")
	}

	userToken, err := localAuth.Get(loginTypeApp)
	if err != nil {
		t.Fatalf("get token by claims %v", err)
	}

	if token != userToken {
		t.Errorf("get token by claims token want %s but get '%s'", token, userToken)
	}
	if token == token3232 {
		t.Errorf("get token by claims token not want %s but get '%s'", token3232, token)
	}

}
func TestLocalGetMultiClaimses(t *testing.T) {
	defer localAuth.CleanCache(loginTypeApp.roleType(), loginTypeApp.Id)
	tokenLen := 0
	for i := LoginTypeWeb; i <= LoginTypeWx; i++ {
		wg.Add(1)
		tokenLen++
		loginTypeApp.setLoginType(int(i))
		go func(i LoginType) {
			localAuth.Generate(loginTypeApp)
			wg.Done()
		}(i)
		wg.Wait()
	}
	userTokens, err := localAuth.getUserTokens(loginTypeApp.roleType(), loginTypeApp.Id)
	if err != nil {
		t.Fatal("get custom claimses generate token is empty \n")
	}
	clas, err := localAuth.getMultiClaimses(userTokens)
	if err != nil {
		t.Fatalf("get custom claimses %v", err)
	}

	if len(userTokens) != tokenLen {
		t.Fatalf("get custom claimses want len %d but get %d", tokenLen, len(userTokens))
	}
	if len(clas) != tokenLen {
		t.Fatalf("get custom claimses want len %d but get %d", tokenLen, len(clas))
	}

}
