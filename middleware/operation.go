package middleware

import (
	"bytes"
	"io/ioutil"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/multi"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func OperationRecord() iris.Handler {
	return func(ctx iris.Context) {
		var body []byte
		var err error

		// 上传文件记录日志文件数据太大
		if !strings.Contains(ctx.Path(), "/api/v1/upload") {
			body, err = ctx.GetBody()
			if err != nil {
				g.ZAPLOG.Error("获取请求内容错误", zap.String("错误:", err.Error()))
			} else {
				ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}
		}

		userId := multi.GetUserId(ctx)
		record := Oplog{
			Ip:     ctx.RemoteAddr(),
			Method: ctx.Method(),
			Path:   ctx.Path(),
			Agent:  ctx.Request().UserAgent(),
			Body:   string(body),
			UserID: userId,
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
		record.ErrorMessage = errorMessage
		record.Status = ctx.GetStatusCode()
		record.Latency = latency
		record.Resp = writer.body.String()

		if err := CreateOplog(record); err != nil {
			g.ZAPLOG.Error("生成日志错误", zap.String("错误:", err.Error()))
		}
	}
}

type responseBodyWriter struct {
	context.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// CreateOplog
func CreateOplog(ol Oplog) error {
	err := database.Instance().Model(&Oplog{}).Create(&ol).Error
	if err != nil {
		g.ZAPLOG.Error("生成系统日志错误", zap.String("错误:", err.Error()))
		return err
	}
	return nil
}

type Oplog struct {
	gorm.Model
	Ip           string        `json:"ip" form:"ip" gorm:"column:ip;comment:请求ip"`
	Method       string        `json:"method" form:"method" gorm:"column:method;comment:请求方法"`
	Path         string        `json:"path" form:"path" gorm:"column:path;comment:请求路径"`
	Status       int           `json:"status" form:"status" gorm:"column:status;comment:请求状态"`
	Latency      time.Duration `json:"latency" form:"latency" gorm:"column:latency;comment:延迟"`
	Agent        string        `json:"agent" form:"agent" gorm:"column:agent;comment:代理"`
	ErrorMessage string        `json:"error_message" form:"error_message" gorm:"column:error_message;comment:错误信息"`
	Body         string        `json:"body" form:"body" gorm:"type:longtext;column:body;comment:请求Body"`
	Resp         string        `json:"resp" form:"resp" gorm:"type:longtext;column:resp;comment:响应Body"`
	UserID       uint          `json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户id"`
}
