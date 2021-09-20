package g

// Model
type Model struct {
	Id        uint   `json:"id"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
}

// ReqId
type ReqId struct {
	Id uint `json:"id" param:"id"`
}

// Paginate
type Paginate struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	OrderBy  string `json:"orderBy"`
	Sort     string `json:"sort"`
}

// Response
type Response struct {
	Code int64       `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

// ErrMsg
type ErrMsg struct {
	Code int64  `json:"code"`
	Msg  string `json:"message"`
}

var (
	NoErr         = ErrMsg{2000, "请求成功"}
	NeedInitErr   = ErrMsg{2001, "前往初始化数据库"}
	AuthErr       = ErrMsg{4001, "认证错误"}
	AuthExpireErr = ErrMsg{4002, "token 过期，请刷新token"}
	AuthActionErr = ErrMsg{4003, "权限错误"}
	ParamErr      = ErrMsg{4004, "参数解析失败，请联系管理员"}
	SystemErr     = ErrMsg{5000, "系统错误，请联系管理员"}
	DataEmptyErr  = ErrMsg{5001, "数据为空"}
	TokenCacheErr = ErrMsg{5002, "TOKEN CACHE 错误"}
)
