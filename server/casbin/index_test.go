package casbin

import "testing"

func TestInstance(t *testing.T) {
	t.Run("test casbin instance", func(t *testing.T) {
		casbin := Instance()
		if casbin == nil {
			t.Error("casbin instance is nil")
		}
	})
}
