package validates

type TenancyRequest struct {
	Name string `json:"name" validate:"required,gte=2,lte=50"  comment:"名称"`
}
