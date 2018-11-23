package controllers

import "github.com/kataras/iris/mvc"

type DocController struct{}

var docView = mvc.View{
	Name: "index.html",
}

// Get will return a predefined view with bind data.
//
// `mvc.Result` is just an interface with a `Dispatch` function.
// `mvc.Response` and `mvc.View` are the built'n result type dispatchers
// you can even create custom response dispatchers by
// implementing the `github.com/kataras/iris/hero#Result` interface.
func (c *DocController) Get() mvc.Result {
	return docView
}
