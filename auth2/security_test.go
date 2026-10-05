package auth2

import (
	"errors"
	"net/http"
	stdhttptest "net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var _ TokenValidator = ValidatorFunc(nil)

func TestNewAgentValidatesConfiguration(t *testing.T) {
	previous := AuthAgent
	defer func() {
		AuthAgent = previous
	}()

	if err := NewAgent(nil); !errors.Is(err, ErrAuthConfigInvalid) {
		t.Fatalf("NewAgent(nil) error = %v, want %v", err, ErrAuthConfigInvalid)
	}
	if err := NewAgent(&Config{Type: "jwt"}); !errors.Is(err, ErrHMACSecretRequired) {
		t.Fatalf("NewAgent without JWT secret error = %v, want %v", err, ErrHMACSecretRequired)
	}
	if err := NewAgent(&Config{Type: "unknown"}); !errors.Is(err, ErrAuthTypeUnsupported) {
		t.Fatalf("NewAgent with unknown type error = %v, want %v", err, ErrAuthTypeUnsupported)
	}
	if err := NewAgent(&Config{HmacSecret: jwtTestSecret}); err != nil {
		t.Fatalf("NewAgent with default JWT type error = %v", err)
	}
	if _, ok := AuthAgent.(*JwtAuth); !ok {
		t.Fatalf("AuthAgent type = %T, want *JwtAuth", AuthAgent)
	}
}

func TestJwtRequiresSecret(t *testing.T) {
	auth := NewJwt(nil)
	claims := validSecurityClaims()

	if _, _, err := auth.Generate(claims); !errors.Is(err, ErrHMACSecretRequired) {
		t.Fatalf("Generate error = %v, want %v", err, ErrHMACSecretRequired)
	}
	if _, err := auth.GetClaims("token"); !errors.Is(err, ErrHMACSecretRequired) {
		t.Fatalf("GetClaims error = %v, want %v", err, ErrHMACSecretRequired)
	}

	auth = NewJwt(jwtTestSecret)
	auth.HmacSecret = nil
	if _, _, err := auth.Generate(claims); !errors.Is(err, ErrHMACSecretRequired) {
		t.Fatalf("Generate after clearing secret error = %v, want %v", err, ErrHMACSecretRequired)
	}

	var nilAuth *JwtAuth
	if _, err := nilAuth.GetClaims("token"); !errors.Is(err, ErrAuthNotInitialized) {
		t.Fatalf("nil GetClaims error = %v, want %v", err, ErrAuthNotInitialized)
	}
}

func TestJwtCopiesSecret(t *testing.T) {
	secret := append([]byte(nil), jwtTestSecret...)
	auth := NewJwt(secret)
	secret[0] ^= 0xff

	if string(auth.HmacSecret) != string(jwtTestSecret) {
		t.Fatal("NewJwt retained the caller's mutable secret slice")
	}
}

func TestJwtRejectsUnexpectedSigningMethodWithoutLeakingToken(t *testing.T) {
	claims := validSecurityClaims()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS384, claims).SignedString(jwtTestSecret)
	if err != nil {
		t.Fatalf("sign HS384 token: %v", err)
	}

	_, err = NewJwt(jwtTestSecret).GetClaims(token)
	if err == nil {
		t.Fatal("GetClaims accepted an HS384 token")
	}
	if strings.Contains(err.Error(), token) {
		t.Fatalf("GetClaims error leaked token: %v", err)
	}
}

func TestVerifierDefaultTokenSource(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := stdhttptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = stdhttptest.NewRequest(http.MethodGet, "/?token=query-token", nil)

	verifier := NewVerifier()
	if token := verifier.RequestToken(ctx); token != "" {
		t.Fatalf("default RequestToken = %q, want empty query token", token)
	}

	ctx.Request.Header.Set("Authorization", "Bearer header-token")
	if token := verifier.RequestToken(ctx); token != "header-token" {
		t.Fatalf("header RequestToken = %q, want header-token", token)
	}

	verifier.Extractors = append(verifier.Extractors, FromQuery)
	ctx.Request.Header.Del("Authorization")
	if token := verifier.RequestToken(ctx); token != "query-token" {
		t.Fatalf("explicit query RequestToken = %q, want query-token", token)
	}
}

func TestVerifierUsesConfiguredValidators(t *testing.T) {
	previous := AuthAgent
	local := NewLocal()
	AuthAgent = local
	defer func() {
		AuthAgent = previous
	}()

	claims := validSecurityClaims()
	token, _, err := local.Generate(claims)
	if err != nil {
		t.Fatalf("Generate local token: %v", err)
	}
	defer local.DelCache(token)

	called := false
	verifier := NewVerifier(ValidatorFunc(func(token []byte, err error) error {
		called = true
		return err
	}))

	recorder := stdhttptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = stdhttptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Request.Header.Set("Authorization", "Bearer "+token)
	verifier.Verify()(ctx)

	if !called {
		t.Fatal("validator configured through NewVerifier was not called")
	}
	if ctx.IsAborted() {
		t.Fatalf("verification aborted: %v", ctx.Errors)
	}
}

func TestVerifierRejectsFailedValidatorAndInvalidatesContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	validatorErr := errors.New("validator rejected token")
	verifier := NewVerifier(ValidatorFunc(func([]byte, error) error {
		return validatorErr
	}))
	handled := error(nil)
	verifier.ErrorHandler = func(ctx *gin.Context, err error) {
		handled = err
		ctx.Abort()
	}

	recorder := stdhttptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = stdhttptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Request.Header.Set("Authorization", "Bearer rejected-token")
	ctx.Set(claimsContextKey, validSecurityClaims())
	ctx.Set(verifiedTokenContextKey, []byte("old-token"))

	verifier.Verify()(ctx)

	if !errors.Is(handled, validatorErr) {
		t.Fatalf("handled error = %v, want %v", handled, validatorErr)
	}
	if got := GetVerifiedToken(ctx); got != nil {
		t.Fatalf("verified token after rejection = %q, want nil", got)
	}
	if got := Get(ctx); got != nil {
		t.Fatalf("claims after rejection = %#v, want nil", got)
	}
}

func TestVerifyTokenRequiresInitializedAgent(t *testing.T) {
	previous := AuthAgent
	AuthAgent = nil
	defer func() {
		AuthAgent = previous
	}()

	_, _, err := NewVerifier().VerifyToken([]byte("token"))
	if !errors.Is(err, ErrAuthNotInitialized) {
		t.Fatalf("VerifyToken error = %v, want %v", err, ErrAuthNotInitialized)
	}
}

func TestIsRoleFailsClosedWithoutTokenOrAgent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := stdhttptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	previous := AuthAgent
	defer func() {
		AuthAgent = previous
	}()

	AuthAgent = NewLocal()
	if IsRole(ctx, RoleAdmin) {
		t.Fatal("IsRole accepted a context without a verified token")
	}

	ctx.Set(verifiedTokenContextKey, []byte("token"))
	AuthAgent = nil
	if IsRole(ctx, RoleAdmin) {
		t.Fatal("IsRole accepted a token without an initialized agent")
	}
}

func TestValidatorFuncCompatibilityMethods(t *testing.T) {
	calls := 0
	validator := ValidatorFunc(func([]byte, error) error {
		calls++
		return nil
	})

	if err := validator.Validater([]byte("token"), nil); err != nil {
		t.Fatalf("Validater error: %v", err)
	}
	if err := validator.ValidateToken([]byte("token"), nil); err != nil {
		t.Fatalf("ValidateToken error: %v", err)
	}
	if calls != 2 {
		t.Fatalf("validator calls = %d, want 2", calls)
	}
}

func validSecurityClaims() *Claims {
	return NewClaims(&Agent{
		Id:        1,
		Username:  "security-test",
		AuthIds:   []string{"1"},
		RoleType:  RoleAdmin,
		LoginType: LoginTypeWeb,
		AuthType:  AuthPwd,
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
	})
}
