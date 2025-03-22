package g

import (
	"os"
	"strings"
)

var (
	// redis mysql mongo cache default config
	TestRedisAddrKey = "IRIS_ADMIN_REDIS_ADDR"
	TestRedisPwdKey  = "IRIS_ADMIN_REDIS_PWD"
	TestMysqlAddrKey = "IRIS_ADMIN_MYSQL_ADDR"
	TestMysqlPwdKey  = "IRIS_ADMIN_MYSQL_PWD"
	TestMysqlNameKey = "IRIS_ADMIN_MYSQL_NAME"
	TestMongoAddrKey = "IRIS_ADMIN_MONGO_ADDR"
	TestDbTypeKey    = "IRIS_ADMIN_DB_TYPE"
	TestRedisAddr    = strings.TrimSpace(os.Getenv(TestRedisAddrKey))
	TestRedisPwd     = strings.TrimSpace(os.Getenv(TestRedisPwdKey))
	TestMysqlAddr    = strings.TrimSpace(os.Getenv(TestMysqlAddrKey))
	TestMysqlPwd     = strings.TrimSpace(os.Getenv(TestMysqlPwdKey))
	TestMysqlName    = strings.TrimSpace(os.Getenv(TestMysqlNameKey))
	TestMongoAddr    = strings.TrimSpace(os.Getenv(TestMongoAddrKey))
	TestDbType       = strings.TrimSpace(os.Getenv(TestDbTypeKey))
)
