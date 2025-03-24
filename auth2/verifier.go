package auth2

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	claimsContextKey        = "gin.auth2.claims"
	verifiedTokenContextKey = "gin.auth2.token"
)

// Get returns the claims decoded by a verifier.
func Get(ctx *gin.Context) *MultiClaims {
	v, b := ctx.Get(claimsContextKey)
	if !b {
		return nil
	}
	tok, ok := v.(*MultiClaims)
	if !ok {
		return nil
	}
	return tok
}

// GetAuthorityType 角色类型
func GetAuthorityType(ctx *gin.Context) int {
	if v := Get(ctx); v != nil {
		return v.AuthorityType
	}
	return 0
}

// GetAuthorityId 角色id
func GetAuthorityId(ctx *gin.Context) []string {
	if v := Get(ctx); v != nil {
		return strings.Split(v.AuthorityId, AuthorityTypeSplit)
	}
	return nil
}

// GetUserId 用户id
func GetUserId(ctx *gin.Context) uint {
	v := Get(ctx)
	if v == nil {
		return 0
	}
	id, err := strconv.Atoi(v.Id)
	if err != nil {
		return 0
	}
	return uint(id)
}

// IsSuperAdmin 用户id
func IsSuperAdmin(ctx *gin.Context) bool {
	v := Get(ctx)
	if v == nil {
		return false
	}
	return v.SuperAdmin
}

// GetUsername 用户名
func GetUsername(ctx *gin.Context) string {
	if v := Get(ctx); v != nil {
		return v.Username
	}
	return ""
}

// GetCreationDate 登录时间
func GetCreationDate(ctx *gin.Context) int64 {
	if v := Get(ctx); v != nil {
		return v.CreationTime
	}
	return 0
}

// GetExpiresIn 有效期
func GetExpiresIn(ctx *gin.Context) int64 {
	if v := Get(ctx); v != nil {
		return v.ExpiresAt
	}
	return 0
}

func GetVerifiedToken(ctx *gin.Context) []byte {
	v, b := ctx.Get(verifiedTokenContextKey)
	if !b {
		return nil
	}
	if tok, ok := v.([]byte); ok {
		return tok
	}
	return nil
}

func IsRole(ctx *gin.Context, authorityType int) bool {
	v := GetVerifiedToken(ctx)
	if v == nil {
		return false
	}
	b, err := AuthDriver.IsRole(string(v), authorityType)
	if err != nil {
		return false
	}
	return b
}

func IsAdmin(ctx *gin.Context) bool {
	return IsRole(ctx, AdminAuthority)
}

type Verifier struct {
	Extractors   []TokenExtractor
	Validators   []TokenValidator
	ErrorHandler func(ctx *gin.Context, err error)
}

func NewVerifier(validators ...TokenValidator) *Verifier {
	return &Verifier{
		Extractors: []TokenExtractor{FromHeader, FromQuery},
		ErrorHandler: func(ctx *gin.Context, err error) {
			ctx.AbortWithError(http.StatusUnauthorized, err)
		},
		Validators: validators,
	}
}

// Invalidate
func (v *Verifier) invalidate(ctx *gin.Context) {
	if verifiedToken := GetVerifiedToken(ctx); verifiedToken != nil {
		ctx.Set(claimsContextKey, "")
		ctx.Set(verifiedTokenContextKey, "")
	}
}

// RequestToken extracts the token from the
func (v *Verifier) RequestToken(ctx *gin.Context) (token string) {
	for _, extract := range v.Extractors {
		if token = extract(ctx); token != "" {
			break // ok we found it.
		}
	}
	return
}

func (v *Verifier) VerifyToken(token []byte, validators ...TokenValidator) ([]byte, *MultiClaims, error) {
	if len(token) == 0 {
		return nil, nil, ErrEmptyToken
	}
	var err error
	for _, validator := range validators {
		// A token validator can skip the builtin validation and return a nil error,
		// in that case the previous error is skipped.
		if err = validator.ValidateToken(token, err); err != nil {
			break
		}
	}
	if err != nil {
		// Exit on parsing standard claims error(when Plain is missing) or standard claims validation error or custom validators.
		return nil, nil, err
	}
	rcc, err := AuthDriver.GetMultiClaims(string(token))
	if err != nil {
		return nil, nil, err
	}
	err = rcc.Valid()
	if err != nil {
		return nil, nil, err
	}
	return token, rcc, nil
}

func (v *Verifier) Verify(validators ...TokenValidator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := []byte(v.RequestToken(ctx))
		verifiedToken, rcc, err := v.VerifyToken(token, validators...)
		if err != nil {
			v.invalidate(ctx)
			v.ErrorHandler(ctx, err)
			return
		}
		ctx.Set(claimsContextKey, rcc)
		ctx.Set(verifiedTokenContextKey, verifiedToken)
		ctx.Next()
	}
}
