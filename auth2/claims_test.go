package auth2

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	cla := New(&Multi{
		Id:            uint(8457585),
		SuperAdmin:    true,
		Username:      "username",
		AuthorityIds:  []string{"999"},
		AuthorityType: AdminAuthority,
		LoginType:     LoginTypeWeb,
		AuthType:      LoginTypeWeb,
		ExpiresAt:     time.Now().Local().Add(RedisSessionTimeoutWeb).Unix(),
	})
	if cla == nil {
		t.Error("claims init return is nil")
	}
}

func TestValid(t *testing.T) {
	cla := New(&Multi{
		Id:            uint(8457585),
		SuperAdmin:    true,
		Username:      "username",
		AuthorityIds:  []string{"999"},
		AuthorityType: AdminAuthority,
		LoginType:     LoginTypeWeb,
		AuthType:      LoginTypeWeb,
		ExpiresAt:     time.Now().Local().Add(RedisSessionTimeoutWeb).Unix(),
	})
	if err := cla.Valid(); err != nil {
		t.Error(err)
	}
}
