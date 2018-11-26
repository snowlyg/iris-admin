package models

import "IrisYouQiKangApi/system"

type AdminUserTranform struct {
	Id               uint   `json:"id"`
	Name             string `json:"name"`
	Username         string `json:"username"`
	IsClient         bool   `json:"is_client"`
	IsFrozen         bool   `json:"is_frozen"`
	IsAudit          bool   `json:"is_audit"`
	IsClientAdmin    bool   `json:"is_client_admin"`
	WechatName       string `json:"wechat_name"`
	WechatAvatar     string `json:"wechat_avatar"`
	Email            string `json:"email"`
	OpenId           string `json:"open_id"`
	WechatVerfiyTime string `json:"wechat_verfiy_time"`
	IsWechatVerfiy   bool   `json:"is_wechat_verfiy"`
	Phone            string `json:"phone"`
	RoleId           uint   `json:"role_id"`
	RoleName         string `json:"role_name"`
	RememberToken    string `json:"remember_token"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	DeletedAt        string `json:"deleted_at"`
}

func TransFormUsers(users []Users) (auts []AdminUserTranform) {

	auts = make([]AdminUserTranform, 0)

	for _, v := range users {
		aut := AdminUserTranform{}
		aut.Id = v.ID
		aut.Name = v.Name
		aut.Username = v.Username

		if v.IsClient == 0 {
			aut.IsClient = false
		} else {
			aut.IsClient = true
		}

		if v.IsFrozen == 0 {
			aut.IsFrozen = false
		} else {
			aut.IsFrozen = true
		}

		if v.IsAudit == 0 {
			aut.IsAudit = false
		} else {
			aut.IsAudit = true
		}

		if v.IsClientAdmin == 0 {
			aut.IsClientAdmin = false
		} else {
			aut.IsClientAdmin = true
		}
		if v.IsWechatVerfiy == 0 {
			aut.IsWechatVerfiy = false
		} else {
			aut.IsWechatVerfiy = true
		}

		aut.WechatName = v.WechatName
		aut.WechatAvatar = v.WechatAvatar
		aut.Email = v.Email
		aut.OpenId = v.OpenId
		aut.Phone = v.Phone
		aut.RoleId = v.Role.ID
		aut.RoleName = v.Role.Name
		aut.RememberToken = v.RememberToken
		aut.CreatedAt = system.Tools.TimeFormat(&v.CreatedAt)
		aut.UpdatedAt = system.Tools.TimeFormat(&v.UpdatedAt)
		if v.DeletedAt == nil {
			aut.DeletedAt = ""
		} else {
			aut.DeletedAt = system.Tools.TimeFormat(v.DeletedAt)
		}
		auts = append(auts, aut)
	}

	return
}
