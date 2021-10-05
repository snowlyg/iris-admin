package zap_server

import (
	"testing"

	"github.com/snowlyg/helper/dir"
	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	t.Run("test zap_server init", func(t *testing.T) {
		Init()
		if !dir.IsExist(CONFIG.Director) {
			t.Error("config dir not exist")
		}
		if !ZAPLOG.Core().Enabled(zap.InfoLevel) {
			t.Error("zaplog not enabled info level")
		}
	})
}
func TestGetEncoderConfig(t *testing.T) {
	t.Run("test zap_server get encoder config", func(t *testing.T) {
		config := getEncoderConfig()
		if config.StacktraceKey != CONFIG.StacktraceKey {
			t.Errorf("zaplog config stacktracekey want %s but get %s", CONFIG.StacktraceKey, config.StacktraceKey)
		}
	})
}
