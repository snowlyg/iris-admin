package file

import (
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/g"
	"go.uber.org/zap"
)

func Upload(ctx iris.Context) {
	f, fh, err := ctx.FormFile("file")
	if err != nil {
		g.ZAPLOG.Error("文件上传失败", zap.String("ctx.FormFile(\"file\")", err.Error()))
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	defer f.Close()

	data, err := UploadFile(ctx, fh)
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: data, Msg: g.NoErr.Msg})
}
