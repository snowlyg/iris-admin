package zap

import (
	"os"
	"path"
	"time"

	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/snowlyg/iris-admin/server/config"
	"go.uber.org/zap/zapcore"
)

// GetWriteSyncer zap logger中加入file-rotatelogs
func GetWriteSyncer() (zapcore.WriteSyncer, error) {
	fileWriter, err := zaprotatelogs.New(
		path.Join(config.CONFIG.Zap.Director, "%Y-%m-%d.log"),
		zaprotatelogs.WithLinkName(config.CONFIG.Zap.LinkName),
		zaprotatelogs.WithMaxAge(7*24*time.Hour),
		zaprotatelogs.WithRotationTime(24*time.Hour),
	)
	if config.CONFIG.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}
