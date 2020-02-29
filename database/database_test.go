package database

import (
	"testing"

	"github.com/snowlyg/IrisAdminApi/config"
)

func TestGetEnforcer(t *testing.T) {
	config.SetAppDriverType("Sqlite")
	if got := GetEnforcer(); got == nil {
		t.Errorf("GetEnforcer() = %v, want %v", got, nil)
	}
}

func TestGetGdb(t *testing.T) {
	config.SetAppDriverType("Sqlite")
	if got := GetGdb(); got == nil {
		t.Errorf("GetGdb() = %v, want %v", got, nil)
	}
}
