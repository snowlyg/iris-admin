package cache

import (
	"testing"
)

func TestInitRedisCluster(t *testing.T) {
	t.Run("init redis cluster", func(t *testing.T) {
		InitRedisCluster([]string{"localhost:6379"}, "foobared")
		redisCluster := GetRedisClusterClient()
		if redisCluster == nil {
			t.Errorf("TestInitRedisCluster error")
		}
	})

}
