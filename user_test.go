package main

import (
	"github.com/kataras/iris/httptest"
	"testing"
)

func TestUsers(t *testing.T) {
	app := NewApp()
	e := httptest.New(t, app)

	/*
		// redirects to /admin without basic auth
		e.GET("/").Expect().Status(httptest.StatusUnauthorized)
		// without basic auth
		e.GET("/admin").Expect().Status(httptest.StatusUnauthorized)

		// with valid basic auth
		e.GET("/admin").WithBasicAuth("myusername", "mypassword").Expect().
			Status(httptest.StatusOK).Body().Equal("/admin myusername:mypassword")
		e.GET("/admin/profile").WithBasicAuth("myusername", "mypassword").Expect().
			Status(httptest.StatusOK).Body().Equal("/admin/profile myusername:mypassword")
		e.GET("/admin/settings").WithBasicAuth("myusername", "mypassword").Expect().
			Status(httptest.StatusOK).Body().Equal("/admin/settings myusername:mypassword")

		// with invalid basic auth
		e.GET("/admin/settings").WithBasicAuth("invalidusername", "invalidpassword").
			Expect().Status(httptest.StatusUnauthorized)
	*/

	login_user_info := map[string]interface{}{
		"username": "admin",
		"password": "admin123456",
	}

	login_response_info := map[string]interface{}{
		"status": true,
		"msg":    "",
	}

	e.POST("/v1/admin/login").WithJSON(login_user_info).
		Expect().Status(httptest.StatusOK).JSON().Object().Equal(login_response_info)

	e.GET("/v1/admin/users/").WithHeader("Bear ", "").Expect().Status(httptest.StatusUnauthorized)
}
