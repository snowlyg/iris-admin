package user

import (
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/validate"
	"github.com/snowlyg/multi"
	"go.uber.org/zap"
)

func Profile(ctx iris.Context) {
	user, err := FindById(database.Instance(), multi.GetUserId(ctx))
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: g.SystemErr.Msg})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: user, Msg: g.NoErr.Msg})
}

func GetUser(ctx iris.Context) {
	var req g.ReqId
	if err := ctx.ReadParams(&req); err != nil {
		g.ZAPLOG.Error("参数解析失败", zap.String("错误:", err.Error()))
		ctx.JSON(g.Response{Code: g.ParamErr.Code, Data: nil, Msg: g.ParamErr.Msg})
		return
	}
	user, err := FindById(database.Instance(), req.Id)
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: g.SystemErr.Msg})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: user, Msg: g.NoErr.Msg})
}

func CreateUser(ctx iris.Context) {
	req := Request{}
	if err := ctx.ReadJSON(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			g.ZAPLOG.Error("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}
	id, err := Create(database.Instance(), req)
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: iris.Map{"id": id}, Msg: g.NoErr.Msg})
}

func UpdateUser(ctx iris.Context) {
	var reqId g.ReqId
	if err := ctx.ReadParams(&reqId); err != nil {
		g.ZAPLOG.Error("参数解析失败", zap.String("错误:", err.Error()))
		ctx.JSON(g.Response{Code: g.ParamErr.Code, Data: nil, Msg: g.ParamErr.Msg})
		return
	}

	var req Request
	if err := ctx.ReadJSON(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			g.ZAPLOG.Error("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}

	err := Update(database.Instance(), reqId.Id, req)
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: nil, Msg: g.NoErr.Msg})
}

func DeleteUser(ctx iris.Context) {
	var req g.ReqId
	if err := ctx.ReadParams(&req); err != nil {
		g.ZAPLOG.Error("参数解析失败", zap.String("错误:", err.Error()))
		ctx.JSON(g.Response{Code: g.ParamErr.Code, Data: nil, Msg: g.ParamErr.Msg})
		return
	}
	err := DeleteById(database.Instance(), req.Id)
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}

	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: nil, Msg: g.NoErr.Msg})
}

func GetAllUsers(ctx iris.Context) {
	var req ReqPaginate
	if err := ctx.ReadQuery(&req); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			g.ZAPLOG.Error("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}

	list, err := Paginate(database.Instance(), req)
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: list, Msg: g.NoErr.Msg})
}

func Logout(ctx iris.Context) {
	token := multi.GetVerifiedToken(ctx)
	if token == nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: "授权凭证为空"})
		return
	}
	err := DelToken(string(token))
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: nil, Msg: g.NoErr.Msg})
}

func Clear(ctx iris.Context) {
	token := multi.GetVerifiedToken(ctx)
	if token == nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: "授权凭证为空"})
		return
	}
	if err := CleanToken(multi.AdminAuthority, string(token)); err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: nil, Msg: g.NoErr.Msg})
}

func ChangeAvatar(ctx iris.Context) {
	avatar := &Avatar{}
	if err := ctx.ReadJSON(avatar); err != nil {
		errs := validate.ValidRequest(err)
		if len(errs) > 0 {
			g.ZAPLOG.Error("参数验证失败", zap.String("错误", strings.Join(errs, ";")))
			ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: strings.Join(errs, ";")})
			return
		}
	}
	err := UpdateAvatar(database.Instance(), multi.GetUserId(ctx), avatar.Avatar)
	if err != nil {
		ctx.JSON(g.Response{Code: g.SystemErr.Code, Data: nil, Msg: err.Error()})
		return
	}
	ctx.JSON(g.Response{Code: g.NoErr.Code, Data: nil, Msg: g.NoErr.Msg})
}
