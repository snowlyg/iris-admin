package controllers

type Response struct {
	Status bool        `json:"status"`
	Msg    interface{} `json:"msg"`
	Data   interface{} `json:"data"`
}

type Lists struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
}

func ApiResource(status bool, objects interface{}, msg string) (apijson *Response) {
	apijson = &Response{Status: status, Data: objects, Msg: msg}
	return
}
