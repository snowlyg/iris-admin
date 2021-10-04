package middleware

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/multi"
)

/**
 * 验证 jwt
 * @method JwtHandler
 */
func JwtHandler() iris.Handler {
	verifier := multi.NewVerifier()
	verifier.Extractors = []multi.TokenExtractor{multi.FromHeader} // extract token only from Authorization: Bearer $token
	verifier.ErrorHandler = func(ctx *context.Context, err error) {
		ctx.JSON(orm.Response{Code: orm.AuthErr.Code, Data: nil, Msg: orm.AuthErr.Msg})
		ctx.StopWithError(http.StatusUnauthorized, err)
	} // extract token only from Authorization: Bearer $token
	return verifier.Verify()
}
