package files

import (
	"path/filepath"
	"testing"
)

func TestGetAbsPath(t *testing.T) {
	absPath := "C:\\Users\\Administrator\\go\\src\\IrisAdminApi\\"

	conf := "config/conf.tml"
	abs := filepath.Join(absPath, conf)
	t.Run("other files", func(t *testing.T) {
		if got := GetAbsPath(conf); got != abs {
			t.Errorf("GetAbsPath() = %v, want %v", got, abs)
		}
	})

	file := "files/file.go"
	abs = filepath.Join(absPath, file)
	t.Run("in files", func(t *testing.T) {
		if got := GetAbsPath(file); got != abs {
			t.Errorf("GetAbsPath() = %v, want %v", got, abs)
		}
	})
}
