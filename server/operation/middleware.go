package operation

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"github.com/snowlyg/multi"
	"go.uber.org/zap"
)

// OperationRecord 操作日志中间件
func OperationRecord() iris.Handler {
	return func(ctx iris.Context) {
		if CONFIG.IsExcept(ctx.Path(), ctx.Method()) && !CONFIG.IsInclude(ctx.Path(), ctx.Method()) {
			ctx.Next()
			return
		}
		var body []byte
		var err error

		body, err = ctx.GetBody()
		if err != nil {
			ctx.Request().Body = ioutil.NopCloser(bytes.NewBufferString(err.Error()))
			zap_server.ZAPLOG.Error("获取请求内容错误", zap.String("错误:", err.Error()))
		} else {
			ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}

		writer := responseBodyWriter{
			ResponseWriter: ctx.ResponseWriter().Clone(),
			body:           &bytes.Buffer{},
		}
		ctx.ResetResponseWriter(writer)
		now := time.Now()

		ctx.Next()

		latency := time.Since(now)
		errorMessage := ""
		if ctx.GetErr() != nil {
			errorMessage = ctx.GetErr().Error()
		}

		record := &Oplog{
			Ip:           ctx.RemoteAddr(),
			Method:       ctx.Method(),
			Path:         ctx.Path(),
			Agent:        ctx.Request().UserAgent(),
			Body:         string(body),
			UserID:       multi.GetUserId(ctx),
			ErrorMessage: errorMessage,
			Status:       ctx.GetStatusCode(),
			Latency:      latency,
		}

		record.Resp = writer.body.String()

		if err := CreateOplog(record); err != nil {
			zap_server.ZAPLOG.Error("生成日志错误", zap.String("CreateOplog()", err.Error()))
		}
	}
}

// responseBodyWriter 响应主体 writer
type responseBodyWriter struct {
	context.ResponseWriter
	body *bytes.Buffer
}

// Write 写入
func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
