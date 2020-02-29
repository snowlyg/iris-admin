package validates

type PermissionRequest struct {
	Name        string `json:"name" validate:"required,gte=4,lte=50" comment:"名称"`
	DisplayName string `json:"display_name" comment:"显示名称"`
	Description string `json:"description" comment:"描述"`
	Act         string `json:"act" comment:"Act"`
}
