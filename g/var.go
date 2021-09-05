package g

import (
	"github.com/go-redis/redis/v8"
	"github.com/snowlyg/iris-admin/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	CONFIG     config.Config
	ZAPLOG     *zap.Logger
	VIPER      *viper.Viper
	CACHE      redis.UniversalClient
	PermRoutes []map[string]string
)
