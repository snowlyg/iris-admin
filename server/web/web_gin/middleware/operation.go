package middleware

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snowlyg/iris-admin/server/operation"
	"github.com/snowlyg/iris-admin/server/zap_server"
	multi "github.com/snowlyg/multi/gin"
	"go.uber.org/zap"
)

func OperationRecord() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// 请求操作不记录
		if ctx.Request.Method == "GET" {
			ctx.Next()
		}

		var body []byte

		var err error
		body, err = ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			zap_server.ZAPLOG.Error("read body from request error:", zap.Any("err", err))
		} else {
			// ioutil.ReadAll 读取数据后重新回写数据
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}

		record := &operation.Oplog{
			Ip:        ctx.ClientIP(),
			Method:    ctx.Request.Method,
			Path:      ctx.Request.URL.Path,
			Agent:     ctx.Request.UserAgent(),
			Body:      string(body),
			UserID:    multi.GetUserId(ctx),
			TenancyId: multi.GetTenancyId(ctx),
		}

		writer := responseBodyWriter{
			ResponseWriter: ctx.Writer,
			body:           &bytes.Buffer{},
		}
		ctx.Writer = writer
		now := time.Now()

		ctx.Next()

		latency := time.Since(now)
		record.ErrorMessage = ctx.Errors.ByType(gin.ErrorTypePrivate).String()
		record.Status = ctx.Writer.Status()
		record.Latency = latency
		record.Resp = writer.body.String()
		if err := operation.CreateOplog(record); err != nil {
			zap_server.ZAPLOG.Error("生成日志错误", zap.String("CreateOplog()", err.Error()))
		}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
