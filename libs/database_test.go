// +build test

package libs

import "testing"

func TestInitDb(t *testing.T) {
	t.Run("TestInitDb", func(t *testing.T) {
		InitDb()
		if Db == nil {
			t.Errorf("TestInitDb error")
		}
	})
}
