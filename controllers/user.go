package controllers

import (
	"IrisYouQiKangApi/logic"
	"IrisYouQiKangApi/models"
	"fmt"
	"github.com/kataras/iris"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

type AdminUserLogin struct {
	Username string `json:"username" validate:"required,gte=4,lte=50"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required,gte=4,lte=50"`
	Phone    string `json:"phone" validate:"required"`
	RoleId   uint   `json:"role_id" validate:"required"`
}

/**
* @api {get} /admin/users/profile 获取登陆用户信息
* @apiName 获取登陆用户信息
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 获取登陆用户信息
* @apiSampleRequest /admin/users/profile
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission 登陆用户
 */
func GetProfile(ctx iris.Context) {
	aun := ctx.Values().Get("auth_user_name")
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(models.ApiJson{Status: true, Data: aun, Msg: "操作成功"})
}

/**
* @api {get} /admin/users/:id 根据id获取用户信息
* @apiName 根据id获取用户信息
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 根据id获取用户信息
* @apiSampleRequest /admin/users/:id
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission 登陆用户
 */
func GetUser(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")

	u := models.Users{}
	u.ID = id

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(u.GetUserById())
}

/**
* @api {post} /admin/login 用户登陆
* @apiName 用户登陆
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 用户登陆
* @apiSampleRequest /admin/login
* @apiParam {string} username 用户名
* @apiParam {string} password 密码
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func UserAdminLogin(ctx iris.Context) {
	aul := new(AdminUserLogin)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(errorData(err))
	} else {
		err1 := validate.Var(aul.Username, "required,min=4,max=20")
		err2 := validate.Var(aul.Password, "required,min=5,max=20")
		if err1 != nil || err2 != nil {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.JSON(errorData(err1, err2))
		} else {
			ctx.StatusCode(iris.StatusOK)
			ctx.JSON(logic.UserAdminCheckLogin(aul.Username, aul.Password))
		}
	}
}

/**
* @api {post} /admin/users/ 新建账号
* @apiName 新建账号
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 新建账号
* @apiSampleRequest /admin/users/
* @apiParam {string} username 用户名
* @apiParam {string} password 密码
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func CreateUser(ctx iris.Context) {
	aul := new(AdminUserLogin)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(errorData(err))
	} else {
		err := validate.Struct(aul)

		if err != nil {

			// This check is only needed when your code could produce
			// an invalid value for validation such as interface with nil
			// value most including myself do not usually have code like this.
			if _, ok := err.(*validator.InvalidValidationError); ok {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.WriteString(err.Error())
				return
			}

			ctx.StatusCode(iris.StatusBadRequest)
			for _, err := range err.(validator.ValidationErrors) {
				fmt.Println()
				fmt.Println(err.Namespace())
				fmt.Println(err.Field())
				fmt.Println(err.StructNamespace()) // Can differ when a custom TagNameFunc is registered or.
				fmt.Println(err.StructField())     // By passing alt name to ReportError like below.
				fmt.Println(err.Tag())
				fmt.Println(err.ActualTag())
				fmt.Println(err.Kind())
				fmt.Println(err.Type())
				fmt.Println(err.Value())
				fmt.Println(err.Param())
				fmt.Println()

				// Or collect these as json objects
				// and send back to the client the collected errors via ctx.JSON
				// {
				// 	"namespace":        err.Namespace(),
				// 	"field":            err.Field(),
				// 	"struct_namespace": err.StructNamespace(),
				// 	"struct_field":     err.StructField(),
				// 	"tag":              err.Tag(),
				// 	"actual_tag":       err.ActualTag(),
				// 	"kind":             err.Kind().String(),
				// 	"type":             err.Type().String(),
				// 	"value":            fmt.Sprintf("%v", err.Value()),
				// 	"param":            err.Param(),
				// }
			}
		} else {
			ctx.StatusCode(iris.StatusOK)
			ctx.JSON(models.CreateUser(aul))
		}

		//err_username := validate.Var(aul.Username, "required,min=4,max=20")
		//err_password := validate.Var(aul.Password, "required,min=5,max=20")
		//err_name := validate.Var(aul.Name, "required,min=5,max=20")
		//err_phone := validate.Var(aul.Phone, "required,min=5,max=20")
		//err_role_id := validate.Var(aul.RoleId, "required,min=5,max=20")

		//if err_username != nil || err_password != nil || err_name != nil || err_phone != nil || err_role_id != nil {
		//	ctx.StatusCode(iris.StatusUnauthorized)
		//	ctx.JSON(errorData(err1, err2))
		//}

	}
}

/**
* @api {post} /admin/users/:id/update 更新账号
* @apiName 更新账号
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 更新账号
* @apiSampleRequest /admin/users/:id/update
* @apiParam {string} username 用户名
* @apiParam {string} password 密码
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func UpdateUser(ctx iris.Context) {
	aul := new(AdminUserLogin)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(errorData(err))
	} else {
		err1 := validate.Var(aul.Username, "required,min=4,max=20")
		err2 := validate.Var(aul.Password, "required,min=5,max=20")
		if err1 != nil || err2 != nil {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.JSON(errorData(err1, err2))
		} else {
			ctx.StatusCode(iris.StatusOK)
			ctx.JSON(logic.UserAdminCheckLogin(aul.Username, aul.Password))
		}
	}
}

/**
* @api {get} /admin/users/:id/frozen 冻结账号
* @apiName 冻结账号
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 冻结账号
* @apiSampleRequest /admin/users/:id/frozen
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func FrozenUser(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")

	u := models.Users{}
	u.ID = id

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(u.FrozenUserById())
}

/**
* @api {get} /admin/users/:id/refrozen 解冻用户
* @apiName 解冻用户
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 解冻用户
* @apiSampleRequest /admin/users/:id/refrozen
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func RefrozenUser(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")

	u := models.Users{}
	u.ID = id

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(u.RefrozenUserById())

}

/**
* @api {get} /admin/users/:id/aduit 设置负责人
* @apiName 设置负责人
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 设置负责人
* @apiSampleRequest /admin/users/:id/aduit
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func SetUserAudit(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")

	u := models.Users{}
	u.ID = id
	u.SetAuditUserById()

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(u.SetAuditUserById())

}

/**
* @api {delete} /admin/users/:id/delete 删除用户
* @apiName 删除用户
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 删除用户
* @apiSampleRequest /admin/users/:id/delete
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func DeleteUser(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")

	u := models.Users{}
	u.ID = id

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(u.DeleteUserById())

}

/**
* @api {get} /users 获取所有的账号
* @apiName 获取所有的账号
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 获取所有的账号
* @apiSampleRequest /users
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func GetAllUsers(ctx iris.Context) {
	cp := Tools.ParseInt(ctx.FormValue("cp"), 1)
	mp := Tools.ParseInt(ctx.FormValue("mp"), 20)
	kw := ctx.FormValue("kw")

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(models.GetAllUsers(kw, cp, mp))
}

/**
* @api {get} /clients 获取所有的客户联系人
* @apiName 获取所有的客户联系人
* @apiGroup Clients
* @apiVersion 1.0.0
* @apiDescription 获取所有的客户联系人
* @apiSampleRequest /clients
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func GetAllClients(ctx iris.Context) {
	cp := Tools.ParseInt(ctx.FormValue("cp"), 1)
	mp := Tools.ParseInt(ctx.FormValue("mp"), 20)
	kw := ctx.FormValue("kw")

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(models.GetAllClients(kw, cp, mp))
}

/**
* @api {get} /logout 用户退出登陆
* @apiName 用户退出登陆
* @apiGroup Users
* @apiVersion 1.0.0
* @apiDescription 用户退出登陆
* @apiSampleRequest /logout
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func UserAdminLogout(ctx iris.Context) {
	json := models.ApiJson{}
	aui := ctx.Values().GetString("auth_user_id")

	uid := uint(Tools.ParseInt(aui, 0))

	json = logic.UserAdminLogout(uid)

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(json)
}
