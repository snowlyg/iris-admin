package models

type AdminPermsTranform struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	GuardName   string `json:"guard_name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at"`
}

func TransFormPerms(perms []Permissions) (auts []AdminPermsTranform) {

	auts = make([]AdminPermsTranform, 0)

	for _, v := range perms {
		aut := AdminPermsTranform{}
		aut.Id = v.ID
		aut.Name = v.Name
		aut.GuardName = v.GuardName
		aut.DisplayName = v.DisplayName
		aut.Description = v.Description
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
