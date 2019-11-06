package controllers

import (
	"fmt"
	"reflect"
)
 
type ApiJson struct {
	Status bool        `json:"status"`
	Msg    interface{} `json:"msg"`
	Data   interface{} `json:"data"`
}

func ApiResource(status bool, objects interface{}, msg string) (apijson *ApiJson) {
	apijson = &ApiJson{Status: status, Data: objects, Msg: msg}
	return
}

func transformer(objects interface{}) interface{} {
	getType := reflect.TypeOf(objects)
	fmt.Println("get Type is :", getType.Name())

	getValue := reflect.ValueOf(objects)
	fmt.Println("get all Fields is:", getValue)

	// 获取方法字段
	// 1. 先获取interface的reflect.Type，然后通过NumField进行遍历
	// 2. 再通过reflect.Type的Field获取其Field
	// 3. 最后通过Field的Interface()得到对应的value
	for i := 0; i < getType.Elem().NumField(); i++ {
		field := getType.Elem().Field(i)
		value := getValue.Elem().Field(i).Interface()
		fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
	}

	// 获取方法
	// 1. 先获取interface的reflect.Type，然后通过.NumMethod进行遍历
	for i := 0; i < getType.NumMethod(); i++ {
		m := getType.Method(i)
		fmt.Printf("%s: %v\n", m.Name, m.Type)
	}

	return objects
}
