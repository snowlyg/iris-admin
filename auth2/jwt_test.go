package auth2

import (
	"fmt"
	"testing"
	"time"
)

var (
	jwtAuth   = NewJwtAuth(nil)
	jwtClaims = New(
		&Multi{
			Id:            uint(8457585),
			Username:      "username",
			AuthorityIds:  []string{"999"},
			AuthorityType: AdminAuthority,
			LoginType:     LoginTypeWeb,
			AuthType:      LoginTypeWeb,
			ExpiresAt:     time.Now().Local().Add(RedisSessionTimeoutWeb).Unix(),
		},
	)
)

func TestJwtGenerateToken(t *testing.T) {
	defer jwtAuth.CleanUserTokenCache(jwtClaims.AuthorityType, jwtClaims.Id)
	t.Run("test generate token", func(t *testing.T) {
		token, _, err := jwtAuth.GenerateToken(jwtClaims)
		if err != nil {
			t.Fatalf("generate token %v", err)
		}
		if token == "" {
			t.Error("generate token is empty")
		}

		t.Logf("token:%s\n", token)

		cc, err := jwtAuth.GetMultiClaims(token)
		if err != nil {
			t.Fatalf("get custom claims  %v", err)
		}

		if cc.Id != jwtClaims.Id {
			t.Errorf("get custom id want %v but get %v", jwtClaims.Id, cc.Id)
		}
		if cc.Username != jwtClaims.Username {
			t.Errorf("get custom username want %v but get %v", jwtClaims.Username, cc.Username)
		}
		if cc.AuthorityId != jwtClaims.AuthorityId {
			t.Errorf("get custom authority_id want %v but get %v", jwtClaims.AuthorityId, cc.AuthorityId)
		}
		if cc.AuthorityType != jwtClaims.AuthorityType {
			t.Errorf("get custom authority_type want %v but get %v", jwtClaims.AuthorityType, cc.AuthorityType)
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
	})
}

func TestJwtDelUserTokenCache(t *testing.T) {
	t.Run("test del user token token", func(t *testing.T) {
		token, _, _ := jwtAuth.GenerateToken(jwtClaims)
		if token == "" {
			t.Error("generate token is empty")
		}
		err := jwtAuth.DelUserTokenCache(token)
		if err != nil {
			t.Errorf("get token by claims token want %v but get %v", nil, err)
		}

	})
}

func TestJwtSetUserTokenMaxCount(t *testing.T) {
	t.Run("test redis set user token max count", func(t *testing.T) {
		err := jwtAuth.SetUserTokenMaxCount(3)
		if err != nil {
			t.Errorf("get token by claims token want %v but get %v", nil, err)
		}
	})
}

func TestJwtGetMultiClaims(t *testing.T) {
	defer jwtAuth.CleanUserTokenCache(jwtClaims.AuthorityType, jwtClaims.Id)
	var token string
	jwtClaims.LoginType = 3
	token, _, err := jwtAuth.GenerateToken(jwtClaims)
	if err != nil {
		t.Fatalf("get custom claims  %v", err)
	}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		jwtClaims.LoginType = i
		go func(i int) {
			jwtAuth.GenerateToken(jwtClaims)
			wg.Done()
		}(i)
		wg.Wait()
	}
	t.Run("test get custom claims", func(t *testing.T) {
		for i := 0; i < 4; i++ {
			go func() {
				cc, err := jwtAuth.GetMultiClaims(token)
				if err != nil {
					t.Errorf("get custom claims  %v", err)
				}
				fmt.Printf("test check token hash get %+v\n", cc)
			}()
		}
		time.Sleep(3 * time.Second)
	})
}

func TestJwtGetTokenByClaims(t *testing.T) {
	t.Run("test get token by claims", func(t *testing.T) {
		_, err := jwtAuth.GetTokenByClaims(jwtClaims)
		if err != nil {
			t.Errorf("get token by claims token want %v but get %v", nil, err)
		}
	})

}
