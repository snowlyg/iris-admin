package casbin

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/helper/dir"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/database"
	"github.com/snowlyg/iris-admin/server/database/orm"
	"github.com/snowlyg/iris-admin/server/zap_server"
	"github.com/snowlyg/multi"
	"go.uber.org/zap"
)

var (
	once     sync.Once
	enforcer *casbin.Enforcer
)

// Instance casbin 单例
func Instance() *casbin.Enforcer {
	once.Do(func() {
		enforcer = GetEnforcer()
	})
	return enforcer
}

// GetEnforcer 获取 casbin.Enforcer
func GetEnforcer() *casbin.Enforcer {
	if database.Instance() == nil {
		zap_server.ZAPLOG.Error(database.ErrDatabaseNotInit.Error())
		return nil
	}
	c, err := gormadapter.NewAdapterByDBUseTableName(database.Instance(), "", "casbin_rule") // Your driver and data source.
	if err != nil {
		zap_server.ZAPLOG.Error("驱动初始化错误", zap.String("gormadapter.NewAdapterByDBUseTableName()", err.Error()))
		return nil
	}

	enforcer, err := casbin.NewEnforcer(filepath.Join(dir.GetCurrentAbPath(), g.CasbinFileName), c)
	if err != nil {
		zap_server.ZAPLOG.Error("初始化失败", zap.String("casbin.NewEnforcer()", err.Error()))
		return nil
	}

	if enforcer == nil {
		zap_server.ZAPLOG.Error("Casbin 未初始化")
		return nil
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		zap_server.ZAPLOG.Error("加载规则失败", zap.String("casbin.LoadPolicy()", err.Error()))
		return nil
	}

	return enforcer
}

// GetRolesForUser 获取角色
func GetRolesForUser(uid uint) []string {
	uids, err := Instance().GetRolesForUser(strconv.FormatUint(uint64(uid), 10))
	if err != nil {
		return []string{}
	}

	return uids
}

// Casbin Casbin 权鉴中间件
func Casbin() iris.Handler {
	return func(ctx *context.Context) {
		check, err := Check(ctx.Request(), strconv.FormatUint(uint64(multi.GetUserId(ctx)), 10))
		if err != nil || !check {
			_, _ = ctx.JSON(orm.Response{Code: orm.AuthActionErr.Code, Data: nil, Msg: err.Error()})
			ctx.StopExecution()
			return
		}

		ctx.Next()
	}
}

// Check checks the username, request's method and path and
// returns true if permission grandted otherwise false.
func Check(r *http.Request, userId string) (bool, error) {
	method := r.Method
	path := r.URL.Path
	ok, err := Instance().Enforce(userId, path, method)
	if err != nil {
		zap_server.ZAPLOG.Error(fmt.Sprintf("验证权限报错：%s-%s-%s", userId, path, method), zap.String("casbinServer.Instance().Enforce()", err.Error()))
		return false, err
	}

	zap_server.ZAPLOG.Debug(fmt.Sprintf("权限：%s-%s-%s", userId, path, method))

	if !ok {
		return ok, errors.New("你未拥有当前操作权限，请联系管理员")
	}
	return ok, nil
}
