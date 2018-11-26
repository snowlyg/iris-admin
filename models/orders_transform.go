package models

import "IrisYouQiKangApi/system"

type AdminOrderTranform struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	StartAt     string `json:"start_at"`
	OrderNumber string `json:"order_number"`
	PlanId      int    `json:"plan_id"`
	CompanyId   int    `json:"company_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at"`
}

func TransFormOrders(orders []Orders) (auts []AdminOrderTranform) {

	auts = make([]AdminOrderTranform, 0)

	for _, v := range orders {
		aut := AdminOrderTranform{}
		aut.Id = v.ID
		aut.Name = v.Name
		aut.Status = v.Status
		aut.OrderNumber = v.OrderNumber
		aut.PlanId = v.PlanId
		aut.CompanyId = v.CompanyId
		aut.CreatedAt = v.CreatedAt.Format("2006-01-02 15:04:05")
		aut.CreatedAt = system.Tools.TimeFormat(&v.CreatedAt)
		aut.UpdatedAt = system.Tools.TimeFormat(&v.UpdatedAt)
		if v.DeletedAt == nil {
			aut.DeletedAt = ""
		} else {
			aut.DeletedAt = system.Tools.TimeFormat(v.DeletedAt)
		}
		if v.StartAt == nil {
			aut.StartAt = ""
		} else {
			aut.StartAt = system.Tools.TimeFormat(v.StartAt)
		}
		auts = append(auts, aut)
	}

	return
}
