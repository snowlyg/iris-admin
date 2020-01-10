package middleware

import (
	"github.com/iris-contrib/middleware/jwt"
)

/**
 * 验证 jwt
 * @method JwtHandler
 */
func JwtHandler() *jwt.Middleware {
	var mySecret = []byte("HS2JDFKhu7Y1av7b")
	j := jwt.New(jwt.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return mySecret, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	return j

}
