package zap_server

import (
	"testing"
)

func TestGetEncoderConfig(t *testing.T) {
	defer Remove()
	t.Run("test zap_server get encoder config", func(t *testing.T) {
		config := getEncoderConfig()
		if config.StacktraceKey != CONFIG.StacktraceKey {
			t.Errorf("zaplog config stacktracekey want %s but get %s", CONFIG.StacktraceKey, config.StacktraceKey)
		}
	})
}
