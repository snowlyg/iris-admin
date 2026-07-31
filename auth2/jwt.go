package auth2

import (
	"fmt"
	"log"

	"github.com/golang-jwt/jwt"
	"github.com/snowlyg/helper/arr"
)

// JwtAuth
type JwtAuth struct {
	HmacSecret []byte
	delToken   arr.ArrayType
	initErr    error
}

// NewJwt returns JWT authentication that fails closed when hmacSecret is empty.
func NewJwt(hmacSecret []byte) *JwtAuth {
	ja := &JwtAuth{
		HmacSecret: append([]byte(nil), hmacSecret...),
		delToken:   arr.NewCheckArrayType(0),
	}
	if len(ja.HmacSecret) == 0 {
		ja.initErr = ErrHMACSecretRequired
	}
	return ja
}

func (ra *JwtAuth) ready() error {
	if ra == nil {
		return ErrAuthNotInitialized
	}
	if ra.initErr != nil {
		return ra.initErr
	}
	if len(ra.HmacSecret) == 0 {
		return ErrHMACSecretRequired
	}
	return nil
}

// Generate
func (ra *JwtAuth) Generate(claims *Claims) (string, int64, error) {
	if err := ra.ready(); err != nil {
		return "", 0, err
	}
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
	if err := ra.ready(); err != nil {
		return nil, err
	}
	if ra.delToken.Check(tokenString) {
		return nil, fmt.Errorf("jwt:token deleted %w", ErrTokenInvalid)
	}
	mc := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("jwt: unexpected signing method %q", token.Method.Alg())
		}
		return ra.HmacSecret, nil
	})
	if err != nil {
		return nil, err
	}
	return mc, nil
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
