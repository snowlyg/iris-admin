package g

import (
	"github.com/go-redis/redis/v8"
	"github.com/snowlyg/iris-admin/config"
	"go.uber.org/zap"
)

var (
	CONFIG config.Config
	ZAPLOG *zap.Logger
	CACHE  redis.UniversalClient
)
