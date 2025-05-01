package auth2

import (
	"fmt"
	"log"

	"github.com/golang-jwt/jwt"
	"github.com/snowlyg/helper/arr"
)

var hmacSampleSecret = []byte("updPA0L2uQ56LwHZoyUX")

// JwtAuth
type JwtAuth struct {
	HmacSecret []byte
	delToken   arr.ArrayType
}

// NewJwt
func NewJwt(hmacSecret []byte) *JwtAuth {
	ja := &JwtAuth{
		HmacSecret: hmacSecret,
		delToken:   arr.NewCheckArrayType(0),
	}
	if ja.HmacSecret == nil {
		ja.HmacSecret = hmacSampleSecret
	}
	return ja
}

// Generate
func (ra *JwtAuth) Generate(claims *Claims) (string, int64, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(ra.HmacSecret)
	if err != nil {
		return "", 0, err
	}
	return tokenString, 0, nil
}

// Token
func (ra *JwtAuth) Token(cla *Claims) (string, error) {
	log.Printf("jwt:get token not support\n")
	return "", nil
}

// GetClaims
func (ra *JwtAuth) GetClaims(tokenString string) (*Claims, error) {
	if ra.delToken.Check(tokenString) {
		return nil, fmt.Errorf("jwt:token deleted %w", ErrTokenInvalid)
	}
	mc := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("incorrect signing method: %v", token.Header["alg"])
		}
		return ra.HmacSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(*Claims); ok && token.Valid {
		return mc, nil
	} else {
		return nil, fmt.Errorf("token[%s]:%w", tokenString, ErrTokenInvalid)
	}
}

// SetLimit
func (ra *JwtAuth) SetLimit(limit int64) error {
	log.Printf("jwt:set max count not support\n")
	return nil
}

// UpdateCacheExpire
func (ra *JwtAuth) UpdateCacheExpire(token string) error {
	log.Printf("jwt:UpdateCacheExpire not support\n")
	return nil
}

// DelCache
func (ra *JwtAuth) DelCache(token string) error {
	ra.delToken.Add(token)
	return nil
}

// CleanCache
func (ra *JwtAuth) CleanCache(roleType RoleType, userId string) error {
	log.Printf("jwt:CleanCache not support")
	return nil
}

// IsRole
func (ra *JwtAuth) IsRole(token string, roleType RoleType) (bool, error) {
	rcc, err := ra.GetClaims(token)
	if err != nil {
		return false, fmt.Errorf("jwt:get User's infomation return error: %w", err)
	}
	return rcc.roleType() == roleType, nil
}

// IsSuperAdmin
func (ra *JwtAuth) IsSuperAdmin(token string) bool {
	rcc, err := ra.GetClaims(token)
	if err != nil {
		log.Printf("jwt:get claims fail:%s\n", err.Error())
		return false
	}
	return rcc.SuperAdmin
}

// Close
func (ra *JwtAuth) Close() {}
