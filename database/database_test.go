package database

import (
	"testing"
)

func TestGetEnforcer(t *testing.T) {
	if got := GetEnforcer();got == nil {
		t.Errorf("GetEnforcer() = %v, want %v", got, nil)
	}
}

func TestGetGdb(t *testing.T) {
	if got := GetGdb();got == nil {
		t.Errorf("GetGdb() = %v, want %v", got, nil)
	}
}
