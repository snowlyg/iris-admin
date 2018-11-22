package models

type AdminPlanTranform struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Editer    string `json:"status"`
	IsParent  bool   `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}

func TransFormPlans(orders []Plans) (auts []AdminPlanTranform) {

	auts = make([]AdminPlanTranform, 0)

	for _, v := range orders {
		aut := AdminPlanTranform{}
		aut.Id = v.ID
		aut.Name = v.Name
		aut.Editer = v.Editer
		if v.IsParent == 0 {
			aut.IsParent = false
		} else {
			aut.IsParent = true
		}
		aut.CreatedAt = Tools.TimeFormat(&v.CreatedAt)
		aut.UpdatedAt = Tools.TimeFormat(&v.UpdatedAt)
		if v.DeletedAt == nil {
			aut.DeletedAt = ""
		} else {
			aut.DeletedAt = Tools.TimeFormat(v.DeletedAt)
		}

		auts = append(auts, aut)
	}

	return
}
