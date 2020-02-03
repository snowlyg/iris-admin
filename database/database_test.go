package database

import (
	"flag"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	SetDatabasePath("rbac_model.conf")

	flag.Parse()
	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestGetEnforcer(t *testing.T) {
	if got := GetEnforcer(); got == nil {
		t.Errorf("GetEnforcer() = %v, want %v", got, nil)
	}
}

func TestGetGdb(t *testing.T) {
	if got := GetGdb(); got == nil {
		t.Errorf("GetGdb() = %v, want %v", got, nil)
	}
}
