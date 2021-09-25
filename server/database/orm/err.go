package orm

import "errors"

var (
	ErrPaginateParam = errors.New("分页查询参数缺失")
	ErrParamValidate = errors.New("参数验证失败")
)
