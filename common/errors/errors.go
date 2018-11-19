package errors

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
)

type Status struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	AccessToken string `json:"access_token"`
}

const (
	//操作成功
	UNMARSHAL_REQBODY_ERR = "unmarshal_reqbody_err"
	//数据库增删改查错误
	DB_CRUD_ERR  = "db_crud_err"
	NEW_USER_ERR = "new_user_err"

	//reids
	REDIS_ERR = "redis_err"

	//jwt
	TOKEN_ERR = "token_err"

	REQUEST_PARAM_ERR = "request_param_err"
)

func NewStatus(status bool, msg string) *Status {
	return &Status{
		Status:  status,
		Message: msg,
	}
}

//把对象转换为json
func (s *Status) ToString() string {
	data, err := json.Marshal(s)
	if err != nil {
		logs.Error("status marshal tostring error:", err)
		return ""
	}
	return string(data)
}

//把对象转换成[]byte类型，用于匹配http返回的类型
func (s *Status) ToBytes() []byte {
	data, err := json.Marshal(s)
	if err != nil {
		logs.Error("status marshal tostring error:", err)
		return []byte("")
	}
	return data
}
