package controllers

import (
	"IrisYouQiKangApi/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"net/http"
	"time"
)

func GetProfile(ctx iris.Context) {

	if (ctx.Method() == "OPTIONS") {
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
	userLogin := new(models.UserLogin)

	if (ctx.Method() == "OPTIONS") {
		return
	}

	if err := ctx.ReadJSON(&userLogin); err != nil {
		//ctx.WriteString(err.Error())
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(iris.Map{
			"status":       false,
			"access_token": "",
			"msg":          err.Error(),
		})

		return
	}

	err1 := validate.Var(userLogin.Username, "required,min=4,max=20")
	err2 := validate.Var(userLogin.Password, "required,min=6,max=20")
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

	user, err := models.UserAdminCheckLogin(userLogin)

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
	claims["userId"] = user.Id
	claims["username"] = user.Username
	claims["name"] = user.Name
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
