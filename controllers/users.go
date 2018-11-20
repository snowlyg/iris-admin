package controllers

import (
	"IrisYouQiKangApi/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"net/http"
	"time"
)

type AdminUserLogin struct {
	Username string
	Password string
}

func GetProfile(ctx iris.Context) {
	if ctx.Method() == "OPTIONS" {
		return
	}

	jwtClaims := ctx.Values().Get("jwt").(*jwt.Token).Claims

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"Code": true,
		"data": jwtClaims,
		"Msg":  "",
	})

	return
}

func UserAdminLogin(ctx iris.Context) {
	aul := new(AdminUserLogin)

	if ctx.Method() == "OPTIONS" {
		return
	}

	if err := ctx.ReadJSON(&aul); err != nil {
		//ctx.WriteString(err.Error())
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(iris.Map{
			"status":       false,
			"access_token": "",
			"msg":          err.Error(),
		})

		return
	}

	err1 := validate.Var(aul.Username, "required,min=4,max=20")
	err2 := validate.Var(aul.Password, "required,min=6,max=20")
	if err1 != nil || err2 != nil {
		//fmt.Println("usernameError:", err1)
		//fmt.Println("passwordError:", err2)
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(iris.Map{
			"status":       false,
			"access_token": "",
			"msg":          errorValidate(),
		})

		return
	}

	u, err := models.UserAdminCheckLogin(aul.Username, aul.Password)

	if err != nil {
		//ctx.WriteString(err.Error())
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(iris.Map{
			"status":       false,
			"access_token": "",
			"msg":          err.Error(),
		})

		return
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = u.Id
	claims["username"] = u.Username
	claims["name"] = u.Name
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		//ctx.WriteString(err.Error())
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(iris.Map{
			"status":       false,
			"access_token": "",
			"msg":          err,
			"expire":       claims["exp"],
		})

		return
	}

	ctx.StatusCode(http.StatusOK)
	ctx.JSON(iris.Map{
		"status":       true,
		"access_token": t,
		"msg":          "",
	})

	return

}

//func GetAllUsers(ctx iris.Context) {
//
//	alldata := models.AllData{
//		Limit: 200,
//	}
//
//	if err := ctx.ReadJSON(&alldata); err != nil {
//		// Handle error.
//	}
//	var (
//		fields []string
//		sortby []string
//		order  []string
//		query  = make(map[string]string)
//	)
//
//	// fields: col1,col2,entity.col3
//	if alldata.Fields != "" {
//		fields = strings.Split(alldata.Fields, ",")
//	}
//
//	// sortby: col1,col2
//	if alldata.Sortby != "" {
//		sortby = strings.Split(alldata.Sortby, ",")
//	}
//	// order: desc,asc
//	if alldata.Order != "" {
//		order = strings.Split(alldata.Order, ",")
//	}
//	// query: k:v,k:v
//	if alldata.Query != "" {
//		for _, cond := range strings.Split(alldata.Query, ",") {
//			kv := strings.SplitN(cond, ":", 2)
//			if len(kv) != 2 {
//				ctx.StatusCode(http.StatusOK)
//				ctx.JSON(iris.Map{
//					"status": false,
//					"data":   "",
//					"msg":    "Error: invalid query key/value pair",
//				})
//
//				return
//			}
//			k, v := kv[0], kv[1]
//			query[k] = v
//		}
//	}
//
//	l, err := models.GetAllUsers(query, fields, sortby, order, alldata.Offset, alldata.Limit)
//	if err != nil {
//		ctx.StatusCode(http.StatusOK)
//		ctx.JSON(iris.Map{
//			"status": true,
//			"data":   "",
//			"msg":    err.Error(),
//		})
//
//		return
//	} else {
//		ctx.StatusCode(http.StatusOK)
//		ctx.JSON(iris.Map{
//			"status": true,
//			"data":   l,
//			"msg":    "",
//		})
//
//		return
//	}
//
//}
//
//func UserAdminLogout(ctx iris.Context) {
//
//}
