package auth2

import (
	"fmt"
	"log"

	"github.com/golang-jwt/jwt"
)

var hmacSampleSecret = []byte("updPA0L2uQ56LwHZoyUX")

// JwtAuth
type JwtAuth struct {
	HmacSecret []byte
}

// NewJwt
func NewJwt(hmacSecret []byte) *JwtAuth {
	ja := &JwtAuth{
		HmacSecret: hmacSecret,
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

// Get 获取用户信息
func (ra *JwtAuth) Get(cla *Claims) (string, error) {
	return "", nil
}

// GetClaims
func (ra *JwtAuth) GetClaims(tokenString string) (*Claims, error) {
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

// SetMaxCount
func (ra *JwtAuth) SetMaxCount(tokenMaxCount int64) error {
	return nil
}

// UpdateCacheExpire
func (ra *JwtAuth) UpdateCacheExpire(token string) error {
	return nil
}

// DelCache
func (ra *JwtAuth) DelCache(token string) error {
	log.Println("auth2: jwt del user token")
	return nil
}

// CleanCache
func (ra *JwtAuth) CleanCache(roleType RoleType, userId string) error {
	return nil
}

// IsRole
func (ra *JwtAuth) IsRole(token string, roleType RoleType) (bool, error) {
	rcc, err := ra.GetClaims(token)
	if err != nil {
		return false, fmt.Errorf("get User's infomation return error: %w", err)
	}
	return rcc.roleType() == roleType, nil
}

// IsSuperAdmin
func (ra *JwtAuth) IsSuperAdmin(token string) bool {
	rcc, err := ra.GetClaims(token)
	if err != nil {
		return false
	}
	return rcc.SuperAdmin
}

// Close
func (ra *JwtAuth) Close() {
}
