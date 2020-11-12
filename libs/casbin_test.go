// +build test

package libs

import "testing"

func TestInitCasbin(t *testing.T) {
	t.Run("TestInitCasbin", func(t *testing.T) {
		InitCasbin()
		if Enforcer == nil {
			t.Errorf("TestInitCasbin error")
		}
	})
}
