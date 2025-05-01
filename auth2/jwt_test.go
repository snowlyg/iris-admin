package auth2

import (
	"testing"
	"time"
)

var (
	jwtAuth   = NewJwt(nil)
	jwtClaims = NewClaims(
		&Agent{
			Id:        uint(8457585),
			Username:  "jwt username",
			AuthIds:   []string{"999"},
			RoleType:  RoleAdmin,
			LoginType: LoginTypeWeb,
			AuthType:  AuthPwd,
			ExpiresAt: time.Now().Local().Add(TimeoutWeb).Unix(),
		},
	)
)

func TestJwtGenerateToken(t *testing.T) {
	defer jwtAuth.CleanCache(jwtClaims.roleType(), jwtClaims.Id)
	token, _, err := jwtAuth.Generate(jwtClaims)
	if err != nil {
		t.Fatalf("generate token %v", err)
	}
	if token == "" {
		t.Error("generate token is empty")
	}

	cc, err := jwtAuth.GetClaims(token)
	if err != nil {
		t.Fatalf("get custom claims fail:%v", err)
	}

	if cc.Id != jwtClaims.Id {
		t.Errorf("get custom id want %v but get %v", jwtClaims.Id, cc.Id)
	}
	if cc.Username != jwtClaims.Username {
		t.Errorf("get custom username want %v but get %v", jwtClaims.Username, cc.Username)
	}
	if cc.AuthId != jwtClaims.AuthId {
		t.Errorf("get custom authority_id want %v but get %v", jwtClaims.AuthId, cc.AuthId)
	}
	if cc.RoleType != jwtClaims.RoleType {
		t.Errorf("get custom authority_type want %v but get %v", jwtClaims.RoleType, cc.RoleType)
	}
	if cc.LoginType != jwtClaims.LoginType {
		t.Errorf("get custom login_type want %v but get %v", jwtClaims.LoginType, cc.LoginType)
	}
	if cc.AuthType != jwtClaims.AuthType {
		t.Errorf("get custom auth_type want %v but get %v", jwtClaims.AuthType, cc.AuthType)
	}
	if cc.CreationTime != jwtClaims.CreationTime {
		t.Errorf("get custom creation_data want %v but get %v", jwtClaims.CreationTime, cc.CreationTime)
	}
	if cc.ExpiresAt != jwtClaims.ExpiresAt {
		t.Errorf("get custom expires_at want %v but get %v", jwtClaims.ExpiresAt, cc.ExpiresAt)
	}
}

func TestJwtSetUserTokenMaxCount(t *testing.T) {
	err := jwtAuth.SetLimit(3)
	if err != nil {
		t.Errorf("get token by claims token want %v but get %v", nil, err)
	}
}

func TestJwtGetMultiClaims(t *testing.T) {
	defer jwtAuth.CleanCache(jwtClaims.roleType(), jwtClaims.Id)
	var token string
	jwtClaims.setLoginType(int(LoginTypeWeb))
	token, _, err := jwtAuth.Generate(jwtClaims)
	if err != nil {
		t.Fatalf("get custom claims  %v", err)
	}
	for i := LoginTypeWeb; i <= LoginTypeDevice; i++ {
		wg.Add(1)
		jwtClaims.setLoginType(int(i))
		go func(i LoginType) {
			jwtAuth.Generate(jwtClaims)
			wg.Done()
		}(i)
		wg.Wait()
	}
	for i := 0; i < 4; i++ {
		go func() {
			_, err := jwtAuth.GetClaims(token)
			if err != nil {
				t.Errorf("get custom claims  %v", err)
			}
		}()
	}
	time.Sleep(3 * time.Second)
}

func TestJwtGetTokenByClaims(t *testing.T) {
	_, err := jwtAuth.Token(jwtClaims)
	if err != nil {
		t.Errorf("get token by claims token want %v but get %v", nil, err)
	}
}

func TestJwtDelUserTokenCache(t *testing.T) {
	token, _, _ := jwtAuth.Generate(jwtClaims)
	if token == "" {
		t.Error("generate token is empty")
	}
	err := jwtAuth.DelCache(token)
	if err != nil {
		t.Errorf("del token fail:%v", err.Error())
	}
	if _, err := jwtAuth.GetClaims(token); err == nil {
		t.Error("del user token fail")
	}
}
