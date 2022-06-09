package casbin

import (
	"strconv"
	"testing"
)

func TestInstance(t *testing.T) {
	t.Run("test casbin instance", func(t *testing.T) {
		casbin := Instance()
		if casbin == nil {
			t.Error("casbin instance is nil")
		}
	})
}
func TestGetEnforcer(t *testing.T) {
	t.Run("test casbin get enforcer", func(t *testing.T) {
		enforcer := GetEnforcer()
		if enforcer == nil {
			t.Error("casbin enforcer is nil")
		}
	})
}

func TestGetRolesForUser(t *testing.T) {
	userId := "888"
	roleId := "2"
	_, err := Instance().AddRoleForUser(userId, roleId)
	if err != nil {
		t.Errorf("add role for user get %v", err.Error())
	}
	userUid, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		t.Errorf("parse uint err %v", err.Error())
	}
	t.Run("test casbin get enforcer", func(t *testing.T) {
		uids := GetRolesForUser(uint(userUid))
		if len(uids) != 1 {
			t.Errorf("get role for user want %+v but get %+v", userId, uids)
		}
		if uids[0] != roleId {
			t.Errorf("get role for user want %s but get %s", roleId, uids[0])
		}
	})
}
