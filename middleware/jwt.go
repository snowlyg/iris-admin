package middleware

import (
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
)

/**
 * 验证 jwt
 * @method JwtHandler
 */
func JwtHandler() *jwtmiddleware.Middleware {
	return jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},

		SigningMethod: jwt.SigningMethodHS256,
	})

}
