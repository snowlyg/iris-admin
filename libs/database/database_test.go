// +build test

package database

import (
	"github.com/snowlyg/blog/libs"
	"testing"
)

func TestInitDb(t *testing.T) {
	t.Run("TestInitDb", func(t *testing.T) {
		libs.InitConfig("")
		if Singleton().Db == nil {
			t.Errorf("TestInitDb error")
		}
		if Singleton().Enforcer == nil {
			t.Errorf("TestInitDb error")
		}
	})
}
