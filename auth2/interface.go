package auth2

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	TokenPrefix      = "GST:"
	BindUserPrefix   = "GSBU:"
	UserPrefix       = "GSU:"
	LimitTokenPrefix = "GT_LIMIT_TOKEN"
)

var (
	AuthTypeSplit           = "-"
	LimitTokenDefault int64 = 10
)

var (
	ErrTokenInvalid        = errors.New("token_invalid")
	ErrEmptyToken          = errors.New("token_empty")
	ErrOverLimit           = errors.New("token_over_limit")
	ErrAuthConfigInvalid   = errors.New("auth_config_invalid")
	ErrAuthTypeUnsupported = errors.New("auth_type_unsupported")
	ErrAuthNotInitialized  = errors.New("auth_not_initialized")
	ErrHMACSecretRequired  = errors.New("hmac_secret_required")
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
	TimeoutWeb    = 4 * time.Hour
	TimeoutApp    = 7 * 24 * time.Hour
	TimeoutWx     = 5 * 52 * 168 * time.Hour
	TimeoutDevice = 5 * 52 * 168 * time.Hour
)

// NewAgent initializes the package authentication agent.
func NewAgent(c *Config) error {
	if c == nil {
		return ErrAuthConfigInvalid
	}
	if c.Max == 0 {
		c.Max = 10
	}
	authType := strings.ToLower(strings.TrimSpace(c.Type))
	if authType == "" {
		authType = "jwt"
	}
	switch authType {
	case "redis":
		agent, err := NewRedis(c.UniversalClient)
		if err != nil {
			return err
		}

		AuthAgent = agent
		err = AuthAgent.SetLimit(c.Max)
		if err != nil {
			return err
		}
	case "local":
		AuthAgent = NewLocal()
		err := AuthAgent.SetLimit(c.Max)
		if err != nil {
			return err
		}
	case "jwt":
		if len(c.HmacSecret) == 0 {
			return ErrHMACSecretRequired
		}
		AuthAgent = NewJwt(c.HmacSecret)
	default:
		return fmt.Errorf("%s: %w", c.Type, ErrAuthTypeUnsupported)
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

// Validater completes the TokenValidator interface.
func (fn ValidatorFunc) Validater(token []byte, err error) error {
	return fn(token, err)
}

// ValidateToken is kept for backward compatibility.
// Deprecated: TokenValidator uses Validater.
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
	Token(claims *Claims) (string, error)
	CleanCache(roleType RoleType, userId string) error
	SetLimit(max int64) error
	IsRole(token string, roleType RoleType) (bool, error)
	IsSuperAdmin(token string) bool
	Close()
}

// getExpire
func getExpire(loginType LoginType) time.Duration {
	switch loginType {
	case LoginTypeWeb:
		return TimeoutWeb
	case LoginTypeWx:
		return TimeoutWx
	case LoginTypeApp:
		return TimeoutApp
	case LoginTypeDevice:
		return TimeoutDevice
	default:
		return TimeoutWeb
	}
}

// getPrefixKey
func getPrefixKey(roleType RoleType, id string) string {
	return fmt.Sprintf("%s%d_%s", UserPrefix, roleType, id)
}
