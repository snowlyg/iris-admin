// +build test

package libs

import "testing"

func TestGetRedisUris(t *testing.T) {
	t.Run("TestGetRedisUris", func(t *testing.T) {
		uris := GetRedisUris()
		if uris != nil && uris[0] == "localhost" {
			t.Errorf("TestInitCasbin error")
		}
	})
}
