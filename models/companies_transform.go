package models

import "IrisYouQiKangApi/system"

type AdminCompanyTranform struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Logo      string `json:"status"`
	Creator   string `json:"status"`
	Preview   bool   `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

func TransFormCompanies(orders []Companies) (auts []AdminCompanyTranform) {

	auts = make([]AdminCompanyTranform, 0)

	for _, v := range orders {
		aut := AdminCompanyTranform{}
		aut.Id = v.ID
		aut.Name = v.Name
		aut.Logo = v.Logo
		aut.Creator = v.Creator
		if v.Preview == 0 {
			aut.Preview = false
		} else {
			aut.Preview = true
		}
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
