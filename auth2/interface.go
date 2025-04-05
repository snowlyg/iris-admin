package auth2

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	GtSessionTokenPrefix        = "GST:"
	GtSessionBindUserPrefix     = "GSBU:"
	GtSessionUserPrefix         = "GSU:"
	GtSessionUserMaxTokenPrefix = "GT_USER_MAX_TOKEN"
)

var (
	AuthTypeSplit                      = "-"
	GtSessionUserMaxTokenDefault int64 = 10
)

var (
	ErrTokenInvalid = errors.New("token_invalid")
	ErrEmptyToken   = errors.New("token_empty")
	ErrOverLimit    = errors.New("token_over_limit")
)

type RoleType int

const (
	RoleNone RoleType = iota
	RoleAdmin
	RoleTenancy
	RoleGeneral
)

type AuthType int

const (
	NoAuth AuthType = iota
	AuthPwd
	AuthCode
	AuthThirdParty
)

type LoginType int

const (
	LoginTypeWeb LoginType = iota
	LoginTypeApp
	LoginTypeWx
	LoginTypeDevice
)

var (
	RedisSessionTimeoutWeb    = 4 * time.Hour
	RedisSessionTimeoutApp    = 7 * 24 * time.Hour
	RedisSessionTimeoutWx     = 5 * 52 * 168 * time.Hour
	RedisSessionTimeoutDevice = 5 * 52 * 168 * time.Hour
)

func NewAgent(c *Config) error {
	if c.Max == 0 {
		c.Max = 10
	}
	switch c.Type {
	case "redis":
		agent, err := NewRedis(c.UniversalClient)
		if err != nil {
			return err
		}

		AuthAgent = agent
		err = AuthAgent.SetMaxCount(c.Max)
		if err != nil {
			return err
		}
	case "local":
		AuthAgent = NewLocal()
		err := AuthAgent.SetMaxCount(c.Max)
		if err != nil {
			return err
		}
	case "jwt":
		AuthAgent = NewJwt(c.HmacSecret)
	default:
		AuthAgent = NewJwt(c.HmacSecret)
	}

	return nil
}

// Agent
type Agent struct {
	Id           uint      `json:"id,omitempty"`
	SuperAdmin   bool      `json:"superAdmin,omitempty"`
	Username     string    `json:"username,omitempty"`
	AuthIds      []string  `json:"authIds,omitempty"`
	RoleType     RoleType  `json:"type,omitempty"`
	LoginType    LoginType `json:"loginType,omitempty"`
	AuthType     AuthType  `json:"authType,omitempty"`
	CreationTime int64     `json:"creationData,omitempty"`
	ExpiresAt    int64     `json:"expiresAt,omitempty"`
}

type Config struct {
	Type            string
	Max             int64
	UniversalClient redis.UniversalClient
	HmacSecret      []byte
}

type (
	// TokenValidator provides further token and claims validation.
	TokenValidator interface {
		// Validater accepts the token, the claims extracted from that
		// and any error that may caused by claims validation (e.g. ErrExpired)
		// or the previous validator.
		// A token validator can skip the builtin validation and return a nil error.
		// Usage:
		//  func(v *myValidator) Validater(token []byte, standardClaims Claims, err error) error {
		//    if err!=nil { return err } <- to respect the previous error
		//    // otherwise return nil or any custom error.
		//  }
		//
		// Look `Blocklist`, `Expected` and `Leeway` for builtin implementations.
		Validater(token []byte, err error) error
	}

	// ValidatorFunc is the interface-as-function shortcut for a TokenValidator.
	ValidatorFunc func(token []byte, err error) error
)

// ValidateToken completes the ValidateToken interface.
// It calls itself.
func (fn ValidatorFunc) ValidateToken(token []byte, err error) error {
	return fn(token, err)
}

var AuthAgent Authentication

// Authentication
type Authentication interface {
	Generate(claims *Claims) (string, int64, error)
	DelCache(token string) error
	UpdateCacheExpire(token string) error
	GetClaims(token string) (*Claims, error)
	Get(claims *Claims) (string, error)
	CleanCache(roleType RoleType, userId string) error
	SetMaxCount(max int64) error
	IsRole(token string, roleType RoleType) (bool, error)
	IsSuperAdmin(token string) bool
	Close()
}

// getExpire
func getExpire(loginType LoginType) time.Duration {
	switch loginType {
	case LoginTypeWeb:
		return RedisSessionTimeoutWeb
	case LoginTypeWx:
		return RedisSessionTimeoutWx
	case LoginTypeApp:
		return RedisSessionTimeoutApp
	case LoginTypeDevice:
		return RedisSessionTimeoutDevice
	default:
		return RedisSessionTimeoutWeb
	}
}

// getPrefixKey
func getPrefixKey(roleType RoleType, id string) string {
	return fmt.Sprintf("%s%d_%s", GtSessionUserPrefix, roleType, id)
}
