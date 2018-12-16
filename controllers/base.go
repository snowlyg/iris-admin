package controllers

type ApiJson struct {
	Status bool        `json:"status"`
	Msg    interface{} `json:"msg"`
	Data   interface{} `json:"data"`
}

func ApiResource(status bool, data interface{}, msg string) (apijson *ApiJson) {
	apijson = &ApiJson{Status: status, Data: data, Msg: msg}
	return
}
