package middleware

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/multi"
)

/**
 * 验证 multi
 * @method MultiHandler
 */
func MultiHandler() iris.Handler {
	verifier := multi.NewVerifier()
	verifier.Extractors = []multi.TokenExtractor{multi.FromHeader} // extract token only from Authorization: Bearer $token
	verifier.ErrorHandler = func(ctx *context.Context, err error) {
		ctx.StopWithError(http.StatusUnauthorized, err)
	}
	return verifier.Verify()
}
