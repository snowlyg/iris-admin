package zap

import (
	"fmt"
	"time"

	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/iris-admin/server/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// level 日志级别
var (
	level  zapcore.Level
	ZAPLOG *zap.Logger
)

// Init 初始化日志服务
func Init() {
	var logger *zap.Logger

	if dir.IsExist(config.CONFIG.Zap.Director) { // 判断是否有Director文件夹
		dir.InsureDir(config.CONFIG.Zap.Director)
	}

	switch config.CONFIG.Zap.Level { // 初始化配置文件的Level
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
	if config.CONFIG.Zap.ShowLine {
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
		StacktraceKey:  config.CONFIG.Zap.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case config.CONFIG.Zap.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		conf.EncodeLevel = zapcore.LowercaseLevelEncoder
	case config.CONFIG.Zap.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		conf.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case config.CONFIG.Zap.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		conf.EncodeLevel = zapcore.CapitalLevelEncoder
	case config.CONFIG.Zap.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		conf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		conf.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return conf
}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	if config.CONFIG.Zap.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore() (core zapcore.Core) {
	writer, err := GetWriteSyncer() // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return
	}
	return zapcore.NewCore(getEncoder(), writer, level)
}

// 自定义日志输出时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(config.CONFIG.Zap.Prefix + "2006/01/02 - 15:04:05.000"))
}

type StringsArray [][]string

// MarshalLogArray 序列化数组日志
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
