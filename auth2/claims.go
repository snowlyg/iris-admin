package auth2

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	ValidationMalformed        uint32 = 1 << iota // Token is malformed
	ValidationUnverifiable                        // Token could not be verified because of signing problems
	ValidationSignatureInvalid                    // Signature validation failed
	ValidationExpired                             // EXP validation failed
	ValidationId
	ValidationUsername
	ValidationAuthId
	ValidationRoleType
	ValidationLoginType
	ValidationAuthType
)

type Claims struct {
	Id           string `json:"id,omitempty" redis:"id"`
	SuperAdmin   bool   `json:"superAdmin,omitempty" redis:"super_admin"`
	Username     string `json:"username,omitempty" redis:"username"`
	AuthId       string `json:"authId,omitempty" redis:"auth_id"`
	RoleType     int    `json:"roleType,omitempty" redis:"role_type"`
	LoginType    int    `json:"loginType,omitempty" redis:"login_type"`
	AuthType     int    `json:"authType,omitempty" redis:"auth_type"`
	CreationTime int64  `json:"creationData,omitempty" redis:"creation_data"`
	ExpiresAt    int64  `json:"expiresAt,omitempty" redis:"expires_at"`
}

func (c *Claims) roleType() RoleType {
	return RoleType(c.RoleType)
}

func (c *Claims) loginType() LoginType {
	return LoginType(c.LoginType)
}

func (c *Claims) authType() AuthType {
	return AuthType(c.AuthType)
}

func (c *Claims) setRoleType(roleType int) {
	c.RoleType = roleType
}

func (c *Claims) setLoginType(loginType int) {
	c.LoginType = loginType
}

func (c *Claims) setAuthType(authType int) {
	c.AuthType = authType
}

func NewClaims(m *Agent) *Claims {
	claims := &Claims{
		Id:           strconv.FormatUint(uint64(m.Id), 10),
		SuperAdmin:   m.SuperAdmin,
		Username:     m.Username,
		AuthId:       strings.Join(m.AuthIds, "-"),
		RoleType:     int(m.RoleType),
		LoginType:    int(m.LoginType),
		AuthType:     int(m.AuthType),
		CreationTime: time.Now().Local().Unix(),
		ExpiresAt:    m.ExpiresAt,
	}
	return claims
}

func (c *Claims) Valid() error {
	vErr := new(jwt.ValidationError)
	now := time.Now().Unix()
	// The claims below are optional, by default, so if they are set to the
	// default value in Go, let's not fail the verification for them.
	if !c.VerifyExpiresAt(now, false) {
		delta := time.Unix(now, 0).Sub(time.Unix(c.ExpiresAt, 0))
		vErr.Inner = fmt.Errorf("claims:token is expired by %v", delta)
		vErr.Errors |= ValidationExpired
	}
	if !c.VerifyId() {
		vErr.Inner = fmt.Errorf("claims:id[%s] is empty", c.Id)
		vErr.Errors |= ValidationId
	}
	if !c.VerifyUsername() {
		vErr.Inner = fmt.Errorf("claims:username[%s] is empty", c.Username)
		vErr.Errors |= ValidationUsername
	}
	if !c.VerifyAuthId() {
		vErr.Inner = fmt.Errorf("claims:authId[%s] is empty", c.AuthId)
		vErr.Errors |= ValidationAuthId
	}
	if !c.VerifyType() {
		vErr.Inner = fmt.Errorf("claims:roleType[%d] is invalid", c.RoleType)
		vErr.Errors |= ValidationRoleType
	}
	if !c.VerifyLoginType() {
		vErr.Inner = fmt.Errorf("claims:loginType[%d] is invalid", c.LoginType)
		vErr.Errors |= ValidationLoginType
	}
	if !c.VerifyAuthType() {
		vErr.Inner = fmt.Errorf("claims:authType[%d] is invalid", c.AuthType)
		vErr.Errors |= ValidationAuthType
	}
	if !valid(vErr) {
		return vErr
	}

	return nil
}

// No errors
func valid(e *jwt.ValidationError) bool {
	return e.Errors == 0
}

// Compares the exp claim against cmp.
// If required is false, this method will return true if the value matches or is unset
func (c *Claims) VerifyExpiresAt(cmp int64, req bool) bool {
	return verifyExp(c.ExpiresAt, cmp, req)
}

func verifyExp(exp int64, now int64, required bool) bool {
	if exp == 0 {
		return !required
	}
	return now <= exp
}

func (c *Claims) VerifyId() bool {
	if id, err := strconv.Atoi(c.Id); err != nil {
		return false
	} else if id > 0 {
		return true
	}
	return false
}

func (c *Claims) VerifyUsername() bool {
	return c.Username != ""
}

func (c *Claims) VerifyAuthId() bool {
	return c.AuthId != ""
}

func (c *Claims) VerifyType() bool {
	return c.RoleType > 0
}

func (c *Claims) VerifyLoginType() bool {
	return LoginType(c.LoginType) >= LoginTypeWeb && LoginType(c.LoginType) <= LoginTypeDevice
}

func (c *Claims) VerifyAuthType() bool {
	return AuthType(c.AuthType) >= NoAuth && AuthType(c.AuthType) <= AuthThirdParty
}
