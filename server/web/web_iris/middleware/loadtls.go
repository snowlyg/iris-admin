package middleware

import (
	"github.com/iris-contrib/middleware/secure"
	"github.com/kataras/iris/v12"
)

// 用https把这个中间件在router里面use一下就好
func LoadTls() iris.Handler {
	middleware := secure.New(secure.Options{
		SSLRedirect: true,
		SSLHost:     "localhost:443",
	})
	return middleware.Handler
}
