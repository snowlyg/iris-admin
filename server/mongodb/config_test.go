package mongodb

import (
	"testing"
)

func TestIsExist(t *testing.T) {
	t.Run("test mongodb config IsExist function", func(t *testing.T) {
		if !IsExist() {
			t.Errorf("config's files is not exist.")
		}
	})
	t.Run("Test GetApplyURI function", func(t *testing.T) {
		want := "mongodb://localhost:27017"
		if applyUrl := CONFIG.GetApplyURI(); applyUrl != want {
			t.Errorf("applyURI want %s but get %s", want, applyUrl)
		}
	})
	t.Run("Test Remove function", func(t *testing.T) {
		if err := Remove(); err != nil {
			t.Error(err)
		}
		if IsExist() {
			t.Errorf("config's files remove is fail.")
		}
	})
}
