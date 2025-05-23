package auth2

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

var testAgent = &Agent{
	Id:         uint(8457585),
	SuperAdmin: true,
	Username:   "username",
	AuthIds:    []string{"999"},
	RoleType:   RoleAdmin,
	LoginType:  LoginTypeWeb,
	AuthType:   AuthPwd,
	ExpiresAt:  time.Now().Local().Add(TimeoutWeb).Unix(),
}

func TestNewClaims(t *testing.T) {
	cla := NewClaims(testAgent)
	if cla == nil {
		t.Fatal("claims init return is nil")
	}
	if cla.Id != "8457585" {
		t.Error("claims id is not 8457585")
	}
	if cla.SuperAdmin != true {
		t.Error("claims super admin is not true")
	}
	if cla.Username != "username" {
		t.Error("claims username is not username")
	}
	if cla.AuthId != "999" {
		t.Error("claims auth ids is not 999")
	}
	if cla.roleType() != RoleAdmin {
		t.Error("claims type is not admin")
	}
	if cla.loginType() != LoginTypeWeb {
		t.Error("claims login type is not web")
	}
	if cla.authType() != AuthPwd {
		t.Error("claims auth type is not web")
	}
	if cla.ExpiresAt != time.Now().Local().Add(TimeoutWeb).Unix() {
		t.Error("claims expires at is not now")
	}

	cla.setAuthType(2)
	if cla.authType() != AuthCode {
		t.Error("claims authType is not authCode")
	}

	cla.setLoginType(1)
	if cla.loginType() != LoginTypeApp {
		t.Error("claims loginType is not loginTypeApp")
	}
	cla.setRoleType(2)
	if cla.roleType() != RoleTenancy {
		t.Error("claims roleType is not roleTenancy")
	}
}

func TestValid(t *testing.T) {
	cla := NewClaims(&Agent{Id: uint(8457585), SuperAdmin: true, Username: "username", AuthIds: []string{"999"}, RoleType: RoleAdmin, LoginType: LoginTypeWeb, AuthType: AuthPwd, ExpiresAt: time.Now().Local().Add(TimeoutWeb).Unix()})
	if err := cla.Valid(); err != nil {
		t.Fatal(err)
	}
	args := []struct {
		agent *Agent
		name  string
		want  uint32
	}{
		{
			name:  "ValidationAuthType",
			agent: &Agent{Id: uint(8457585), SuperAdmin: true, Username: "username", AuthIds: []string{"999"}, AuthType: 99, RoleType: RoleAdmin, LoginType: LoginTypeWeb, ExpiresAt: time.Now().Local().Add(TimeoutWeb).Unix()},
			want:  ValidationAuthType,
		},
		{
			name:  "ValidationLoginType",
			agent: &Agent{Id: uint(8457585), SuperAdmin: true, Username: "username", AuthIds: []string{"999"}, LoginType: 99, RoleType: RoleAdmin, AuthType: AuthPwd, ExpiresAt: time.Now().Local().Add(TimeoutWeb).Unix()},
			want:  ValidationLoginType,
		},
		{
			name:  "ValidationRoleType",
			agent: &Agent{Id: uint(8457585), SuperAdmin: true, Username: "username", AuthIds: []string{"999"}, LoginType: LoginTypeApp, RoleType: -1, AuthType: AuthPwd, ExpiresAt: time.Now().Local().Add(TimeoutWeb).Unix()},
			want:  ValidationRoleType,
		},
		{
			name:  "ValidationAuthId",
			agent: &Agent{Id: uint(8457585), SuperAdmin: true, Username: "username", AuthIds: []string{""}, LoginType: LoginTypeApp, RoleType: RoleAdmin, AuthType: AuthPwd, ExpiresAt: time.Now().Local().Add(TimeoutWeb).Unix()},
			want:  ValidationAuthId,
		},
		{
			name:  "ValidationUsername",
			agent: &Agent{Id: uint(8457585), SuperAdmin: true, Username: "", AuthIds: []string{"1"}, LoginType: LoginTypeApp, RoleType: RoleAdmin, AuthType: AuthPwd, ExpiresAt: time.Now().Local().Add(TimeoutWeb).Unix()},
			want:  ValidationUsername,
		},
		{
			name:  "ValidationId",
			agent: &Agent{Id: 0, SuperAdmin: true, Username: "username", AuthIds: []string{"1"}, LoginType: LoginTypeApp, RoleType: RoleAdmin, AuthType: AuthPwd, ExpiresAt: time.Now().Local().Add(TimeoutWeb).Unix()},
			want:  ValidationId,
		},
	}

	for _, arg := range args {
		t.Run(arg.name, func(t *testing.T) {
			cla = NewClaims(arg.agent)
			if err := cla.Valid(); err == nil {
				t.Fatal("error is nil")
			} else {

				if v, ok := err.(*jwt.ValidationError); !ok {
					t.Fatalf("%s %s", reflect.TypeOf(err).String(), err.Error())
				} else if v.Errors != arg.want {
					t.Fatalf("%d %d %s", v.Errors, arg.want, v.Error())
				}
			}
		})
	}
}
