package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	multi "github.com/snowlyg/multi/v2"
)

// init 初始化认证驱动
// 驱动类型： 可选 redis ,local,jwt
func init() {
	options := &redis.UniversalOptions{
		Addrs:       []string{"127.0.0.1:6379"},
		Password:    "",
		PoolSize:    10,
		IdleTimeout: 300 * time.Second,
		// Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
		// 	conn, err := net.Dial(network, addr)
		// 	if err == nil {
		// 		go func() {
		// 			time.Sleep(5 * time.Second)
		// 			conn.Close()
		// 		}()
		// 	}
		// 	return conn, err
		// },
	}

	err := multi.InitDriver(&multi.Config{
		DriverType:      "redis",
		TokenMaxCount:   10,
		UniversalClient: redis.NewUniversalClient(options)})
	if err != nil {
		panic(fmt.Sprintf("auth is not init get err %v\n", err))
	}
}

func auth() gin.HandlerFunc {
	verifier := multi.NewVerifier()
	verifier.Extractors = []multi.TokenExtractor{multi.FromHeader} // extract token only from Authorization: Bearer $token
	return verifier.Verify()
}

func main() {
	app := gin.New()

	app.GET("/", generateToken())

	protectedAPI := app.Group("/protected")
	// Register the verify middleware to allow access only to authorized clients.
	protectedAPI.Use(auth())
	// ^ or UseRouter(verifyMiddleware) to disallow unauthorized http error handlers too.

	protectedAPI.GET("/", protected)
	// Invalidate the token through server-side, even if it's not expired yet.
	protectedAPI.GET("/logout", logout)

	// http://localhost:8080
	// http://localhost:8080/protected (or Authorization: Bearer $token)
	// http://localhost:8080/protected/logout
	// http://localhost:8080/protected (401)
	app.Run(":8080")
}

func generateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := &multi.MultiClaims{
			Id:            "1",
			Username:      "your name",
			AuthorityId:   "your authority id",
			AuthorityType: multi.AdminAuthority,
			LoginType:     multi.LoginTypeWeb,
			AuthType:      multi.AuthPwd,
			CreationTime:  time.Now().Local().Unix(),
			ExpiresAt:     time.Now().Local().Add(multi.RedisSessionTimeoutWeb).Unix(),
		}

		token, _, err := multi.AuthDriver.GenerateToken(claims)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.String(200, token)
	}
}

func protected(ctx *gin.Context) {
	claims := multi.Get(ctx)
	ctx.JSON(http.StatusOK, fmt.Sprintf("claims=%+v\n", claims))
}

func logout(ctx *gin.Context) {
	token := multi.GetVerifiedToken(ctx)
	if token == nil {
		ctx.String(http.StatusOK, "授权凭证为空")
		return
	}
	err := multi.AuthDriver.DelUserTokenCache(string(token))
	if err != nil {
		ctx.JSON(http.StatusOK, err.Error())
		return
	}
	ctx.String(http.StatusOK, "token invalidated, a new token is required to access the protected API")
}
