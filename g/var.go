package g

import (
	"os"
	"strings"
)

var (
	TestRedisPwd  = strings.TrimSpace(os.Getenv("redisProPwd"))
	TestMysqlAddr = strings.TrimSpace(os.Getenv("mysqlAddr"))
	TestMysqlPwd  = strings.TrimSpace(os.Getenv("mysqlPwd"))
	TestMongoAddr = strings.TrimSpace(os.Getenv("mongoAddr"))
	TestDbType    = strings.TrimSpace(os.Getenv("dbType"))
)
