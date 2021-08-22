package middleware

// import (
// 	"bytes"
// 	"net/http"
// 	"time"

// 	"github.com/kataras/iris/v12"
// 	"github.com/kataras/iris/v12/context"
// 	"github.com/snowlyg/iris-admin/application/libs/logging"
// 	"github.com/snowlyg/iris-admin/application/models"
// 	"github.com/snowlyg/iris-admin/service/dao"
// )

// func OperationRecord() iris.Handler {
// 	return func(ctx iris.Context) {
// 		var body []byte
// 		var userId uint
// 		if ctx.Method() != http.MethodGet {
// 			var err error
// 			body, err = ctx.GetBody()
// 			if err != nil {
// 				logging.ErrorLogger.Errorf("read body from request error：", err)
// 			} else {
// 				ctx.Recorder().SetBody(body)
// 			}
// 		}

// 		userId, err := dao.GetAuthId(ctx)
// 		if err != nil {
// 			logging.ErrorLogger.Errorf("get auth id error:", err)
// 		}

// 		record := models.Oplog{
// 			Ip:     ctx.RemoteAddr(),
// 			Method: ctx.Method(),
// 			Path:   ctx.Path(),
// 			Agent:  ctx.Request().UserAgent(),
// 			Body:   string(body),
// 			UserID: userId,
// 		}

// 		writer := responseBodyWriter{
// 			ResponseWriter: ctx.ResponseWriter().Clone(),
// 			body:           &bytes.Buffer{},
// 		}
// 		ctx.ResetResponseWriter(writer)
// 		now := time.Now()

// 		ctx.Next()

// 		latency := time.Since(now)
// 		errorMessage := ""
// 		if ctx.GetErr() != nil {
// 			errorMessage = ctx.GetErr().Error()
// 		}
// 		record.ErrorMessage = errorMessage
// 		record.Status = ctx.GetStatusCode()
// 		record.Latency = latency
// 		record.Resp = writer.body.String()

// 		if err := dao.CreateOplog(record); err != nil {
// 			logging.ErrorLogger.Errorf("create operation record error:", err)
// 		}
// 	}
// }

// type responseBodyWriter struct {
// 	context.ResponseWriter
// 	body *bytes.Buffer
// }

// func (r responseBodyWriter) Write(b []byte) (int, error) {
// 	r.body.Write(b)
// 	return r.ResponseWriter.Write(b)
// }
