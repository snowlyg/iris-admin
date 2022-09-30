package zap_server

import (
	"fmt"
	"time"

	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/iris-admin/server/viper_server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const ConfigFileName = "zap.yaml"

// level log level
var (
	level  zapcore.Level
	ZAPLOG *zap.Logger
)

// init
func init() {
	var logger *zap.Logger
	viper_server.Init(getViperConfig())

	if !dir.IsExist(CONFIG.Director) {
		dir.InsureDir(CONFIG.Director)
	}
	switch CONFIG.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	if level == zap.DebugLevel || level == zap.ErrorLevel {
		logger = zap.New(getEncoderCore(), zap.AddStacktrace(level))
	} else {
		logger = zap.New(getEncoderCore())
	}
	if CONFIG.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	ZAPLOG = logger
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig() (conf zapcore.EncoderConfig) {
	conf = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  CONFIG.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case CONFIG.EncodeLevel == "LowercaseLevelEncoder":
		conf.EncodeLevel = zapcore.LowercaseLevelEncoder
	case CONFIG.EncodeLevel == "LowercaseColorLevelEncoder":
		conf.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case CONFIG.EncodeLevel == "CapitalLevelEncoder":
		conf.EncodeLevel = zapcore.CapitalLevelEncoder
	case CONFIG.EncodeLevel == "CapitalColorLevelEncoder":
		conf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		conf.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return conf
}

// getEncoder
func getEncoder() zapcore.Encoder {
	if CONFIG.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

// getEncoderCore
func getEncoderCore() (core zapcore.Core) {
	writer, err := GetWriteSyncer()
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return
	}
	return zapcore.NewCore(getEncoder(), writer, level)
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(CONFIG.Prefix + "2006/01/02 - 15:04:05.000"))
}

type StringsArray [][]string

// MarshalLogArray
func (ss StringsArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range ss {
		for ii := range ss[i] {
			arr.AppendString(ss[i][ii])
		}
	}
	return nil
}

// Strings constructs a field that carries a slice of strings.
func Strings(key string, ss [][]string) zap.Field {
	return zap.Array(key, StringsArray(ss))
}
