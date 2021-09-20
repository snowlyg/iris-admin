package g

import (
	"github.com/go-redis/redis/v8"
	"github.com/snowlyg/iris-admin/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	CONFIG     config.Config         // 配置
	ZAPLOG     *zap.Logger           // 日志
	VIPER      *viper.Viper          // viper
	CACHE      redis.UniversalClient // 缓存
	PermRoutes []map[string]string   // 权限路由
)
