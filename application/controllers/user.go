package controllers

//func GetProfile(ctx iris.Context) {
//	ctx.StatusCode(iris.StatusOK)
//	sess := ctx.Values().Get("sess").(*auth.SessionV2)
//	id := uint(libs.ParseInt(sess.UserId, 10))
//	s := &easygorm.Search{
//		Fields: []*easygorm.Field{
//			{
//				Key:       "id",
//				Condition: "=",
//				Value:     id,
//			},
//		},
//	}
//	user, err := models.GetUser(s)
//	if err != nil {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, err.Error()))
//		return
//	}
//	_, _ = ctx.JSON(response.NewResponse(200, user, "操作成功"))
//}
//
//func GetAdminInfo(ctx iris.Context) {
//	ctx.StatusCode(iris.StatusOK)
//	s := &easygorm.Search{
//		Fields: []*easygorm.Field{
//			{
//				Key:       "username",
//				Condition: "=",
//				Value:     "username",
//			},
//		},
//	}
//	user, err := models.GetUser(s)
//	if err != nil {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, err.Error()))
//		return
//	}
//	_, _ = ctx.JSON(response.NewResponse(200, map[string]string{"avatar": user.Avatar}, "操作成功"))
//}
//
//func ChangeAvatar(ctx iris.Context) {
//	ctx.StatusCode(iris.StatusOK)
//	sess := ctx.Values().Get("sess").(*auth.SessionV2)
//	id := uint(libs.ParseInt(sess.UserId, 10))
//
//	avatar := new(models.Avatar)
//	if err := ctx.ReadJSON(avatar); err != nil {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, err.Error()))
//		return
//	}
//
//	err := validates.Validate.Struct(*avatar)
//	if err != nil {
//		errs := err.(validator.ValidationErrors)
//		for _, e := range errs.Translate(validates.ValidateTrans) {
//			if len(e) > 0 {
//				_, _ = ctx.JSON(response.NewResponse(400, nil, e))
//				return
//			}
//		}
//	}
//
//	user := &models.User{}
//	user.ID = id
//	user.Avatar = avatar.Avatar
//	err = models.UpdateUserById(id, user)
//	if err != nil {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, err.Error()))
//		return
//	}
//	_, _ = ctx.JSON(response.NewResponse(200, user, "操作成功"))
//}
//
//func GetUser(ctx iris.Context) {
//	ctx.StatusCode(iris.StatusOK)
//	id, _ := ctx.Params().GetUint("id")
//	s := &easygorm.Search{
//		Fields: []*easygorm.Field{
//			{
//				Key:       "id",
//				Condition: "=",
//				Value:     id,
//			},
//		},
//	}
//	user, err := models.GetUser(s)
//	if err != nil {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, err.Error()))
//		return
//	}
//	_, _ = ctx.JSON(response.NewResponse(200, user, "操作成功"))
//}
//
//func CreateUser(ctx iris.Context) {
//	ctx.StatusCode(iris.StatusOK)
//	user := new(models.User)
//	if err := ctx.ReadJSON(user); err != nil {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, err.Error()))
//		return
//	}
//
//	err := validates.Validate.Struct(*user)
//	if err != nil {
//		errs := err.(validator.ValidationErrors)
//		for _, e := range errs.Translate(validates.ValidateTrans) {
//			if len(e) > 0 {
//				_, _ = ctx.JSON(response.NewResponse(400, nil, e))
//				return
//			}
//		}
//	}
//
//	err = user.CreateUser()
//	if err != nil {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, err.Error()))
//		return
//	}
//
//	if user.ID == 0 {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, "操作失败"))
//		return
//	}
//	_, _ = ctx.JSON(response.NewResponse(200, user, "操作成功"))
//	return
//
//}
//
//func UpdateUser(ctx iris.Context) {
//	ctx.StatusCode(iris.StatusOK)
//	user := new(models.User)
//
//	if err := ctx.ReadJSON(user); err != nil {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, err.Error()))
//	}
//
//	err := validates.Validate.Struct(*user)
//	if err != nil {
//		errs := err.(validator.ValidationErrors)
//		for _, e := range errs.Translate(validates.ValidateTrans) {
//			if len(e) > 0 {
//				_, _ = ctx.JSON(response.NewResponse(400, nil, e))
//				return
//			}
//		}
//	}
//
//	id, _ := ctx.Params().GetUint("id")
//	if user.Username == "username" {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, "不能编辑管理员"))
//		return
//	}
//
//	err = models.UpdateUserById(id, user)
//	if err != nil {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, err.Error()))
//		return
//	}
//	_, _ = ctx.JSON(response.NewResponse(200, user, "操作成功"))
//}
//
//func DeleteUser(ctx iris.Context) {
//	ctx.StatusCode(iris.StatusOK)
//	id, _ := ctx.Params().GetUint("id")
//
//	err := models.DeleteUser(id)
//	if err != nil {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, err.Error()))
//		return
//	}
//	_, _ = ctx.JSON(response.NewResponse(200, nil, "删除成功"))
//}

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
//func GetAllUsers(ctx iris.Context) {
//	ctx.StatusCode(iris.StatusOK)
//	name := ctx.FormValue("name")
//
//	s := libs.GetCommonListSearch(ctx)
//	s.Fields = append(s.Fields, easygorm.GetField("name", name))
//	users, count, err := models.GetAllUsers(s)
//	if err != nil {
//		_, _ = ctx.JSON(response.NewResponse(400, nil, err.Error()))
//		return
//	}
//
//	_, _ = ctx.JSON(response.NewResponse(200, map[string]interface{}{"items": users, "total": count, "limit": s.Limit}, "操作成功"))
//
//}
//
//func usersTransform(users []*models.User) []*transformer.User {
//	var us []*transformer.User
//	for _, user := range users {
//		u := userTransform(user)
//		us = append(us, u)
//	}
//	return us
//}
//
//func userTransform(user *models.User) *transformer.User {
//	u := &transformer.User{}
//	g := gf.NewTransform(u, user, time.RFC3339)
//	_ = g.Transformer()
//
//	roleIds := easygorm.GetRolesForUser(user.ID)
//	var ris []int
//	for _, roleId := range roleIds {
//		ri, _ := strconv.Atoi(roleId)
//		ris = append(ris, ri)
//	}
//	s := &easygorm.Search{
//		Fields: []*easygorm.Field{
//			{
//				Key:       "id",
//				Condition: "IN",
//				Value:     ris,
//			},
//		},
//	}
//	roles, _, err := models.GetAllRoles(s)
//	if err == nil {
//		u.Roles = rolesTransform(roles)
//	}
//	return u
//}
