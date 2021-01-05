package response

type Response struct {
	Code int64       `json:"code"`
	Msg  interface{} `json:"message"`
	Data interface{} `json:"data"`
}

type MoreResponse struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
}

func NewResponse(code int64, objects interface{}, msg string) *Response {
	return &Response{Code: code, Data: objects, Msg: msg}
}

type ErrMsg struct {
	Code int64
	Msg  string
}

var (
	NoErr         = ErrMsg{2000, "请求成功"}
	AuthErr       = ErrMsg{4001, "认证错误"}
	AuthExpireErr = ErrMsg{4002, "token 过期，请刷新token"}
	AuthActionErr = ErrMsg{4003, "权限错误"}
	SystemErr     = ErrMsg{5000, "系统错误，请联系管理员"}
	DataEmptyErr  = ErrMsg{5001, "数据为空"}
	TokenCacheErr = ErrMsg{5002, "TOKEN CACHE 错误"}
)
