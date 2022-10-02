package cache

import (
	"testing"
)

func TestIsExist(t *testing.T) {
	t.Run("Test Remove function", func(t *testing.T) {
		if err := Remove(); err != nil {
			t.Error(err)
		}
		if IsExist() {
			t.Errorf("config's files remove is fail.")
		}
	})
}
