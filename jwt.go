package main

import "github.com/dgrijalva/jwt-go"
import jwtmiddleware "github.com/iris-contrib/middleware/jwt"

func jwtHandler()  *jwtmiddleware.Middleware {
	return jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
		//Extractor: jwtmiddleware.FromParameter("auth_code"),
		//EnableAuthOnOptions : false,
		//CredentialsOptional : true,

		SigningMethod: jwt.SigningMethodHS256,
	})

}
