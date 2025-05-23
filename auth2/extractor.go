package auth2

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// TokenExtractor is a function that takes a context as input and returns
// a token. An empty string should be returned if no token found
// without additional information.
type TokenExtractor func(*gin.Context) string

// FromHeader is a token extractor.
// It reads the token from the Authorization request header of form:
// Authorization: "Bearer {token}".
func FromHeader(ctx *gin.Context) string {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// pure check: authorization header format must be Bearer {token}
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return ""
	}

	return authHeaderParts[1]
}

// FromQuery is a token extractor.
// It reads the token from the "token" url query parameter.
func FromQuery(ctx *gin.Context) string {
	return ctx.Query("token")
}

// FromJSON is a token extractor.
// Reads a json request body and extracts the json based on the given field.
// The request content-type should contain the: application/json header value, otherwise
// this method will not try to read and consume the body.
func FromJSON(jsonKey string) TokenExtractor {
	return func(ctx *gin.Context) string {
		if ctx.ContentType() != binding.MIMEJSON {
			log.Printf("extractor: content-type %s not supported\n", ctx.ContentType())
			return ""
		}

		var m gin.H
		if err := ctx.BindJSON(&m); err != nil {
			log.Println("extractor: bind json error:", err.Error())
			return ""
		}

		if m == nil {
			log.Println("extractor: json is empty")
			return ""
		}

		v, ok := m[jsonKey]
		if !ok {
			log.Printf("extractor: key %s not found\n", jsonKey)
			return ""
		}

		tok, ok := v.(string)
		if !ok {
			log.Printf("extractor: key %s value:[%v] is not a string\n", jsonKey, v)
			return ""
		}
		return tok
	}
}
