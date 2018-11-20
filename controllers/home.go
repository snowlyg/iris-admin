package controllers

import (
	"github.com/kataras/iris"
)

func GetHomeData(ctx iris.Context) {

	//homedate := new(models.HomeDate)
	//
	//var (
	//	fields []string
	//	sortby []string
	//	order  []string
	//	query        = make(map[string]string)
	//	limit  int64 = 10
	//	offset int64
	//)
	//
	//companies, err1 := models.GetAllCompanies(query, fields, sortby, order, offset, limit)
	//orders, err2 := models.GetAllOrders(query, fields, sortby, order, offset, limit)
	//clients, err3 := models.GetAllUsers(query, fields, sortby, order, offset, limit)
	//
	//if err1 != nil {
	//	ctx.StatusCode(http.StatusOK)
	//	ctx.JSON(iris.Map{
	//		"status": false,
	//		"data":   "",
	//		"msg":    err1.Error(),
	//	})
	//
	//	return
	//}
	//
	//if err2 != nil {
	//	ctx.StatusCode(http.StatusOK)
	//	ctx.JSON(iris.Map{
	//		"status": false,
	//		"data":   "",
	//		"msg":    err2.Error(),
	//	})
	//
	//	return
	//}
	//if err3 != nil {
	//	ctx.StatusCode(http.StatusOK)
	//	ctx.JSON(iris.Map{
	//		"status": false,
	//		"data":   "",
	//		"msg":    err2.Error(),
	//	})
	//
	//	return
	//}
	//homedate.Clients = len(clients)
	//homedate.Companies = len(companies)
	//homedate.Orders = len(orders)
	//
	//ctx.StatusCode(http.StatusOK)
	//ctx.JSON(iris.Map{
	//	"status": true,
	//	"data":   homedate,
	//})
	//
	//return

}
