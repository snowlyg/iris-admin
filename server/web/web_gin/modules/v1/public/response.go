package public

type LoginResponse struct {
	Data  interface{} `json:"data"`
	Token string      `json:"AccessToken"`
}

type MiniCodeResponse struct {
	SessionKey string `json:"sessionKey" form:"sessionKey" `
	OpenId     string `json:"openId" form:"openId" `
	UnionId    string `json:"unionId" form:"unionId"`
}
