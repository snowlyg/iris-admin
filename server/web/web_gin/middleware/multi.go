package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/snowlyg/iris-admin/server/web/web_gin/response"
	multi "github.com/snowlyg/multi/gin"
)

func Auth() gin.HandlerFunc {
	verifier := multi.NewVerifier()
	verifier.Extractors = []multi.TokenExtractor{multi.FromHeader, multi.FromQuery} // extract token  from Authorization: Bearer $token and query ?token=
	verifier.ErrorHandler = func(ctx *gin.Context, err error) {
		response.UnauthorizedFailWithMessage(err.Error(), ctx)
		ctx.Abort()
	} // extract token only from Authorization: Bearer $token
	return verifier.Verify()
}
