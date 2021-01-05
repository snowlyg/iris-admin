package dao

type Dao interface {
	ModelName() string
	All(name, sort, orderBy string, page, pageSize int) (map[string]interface{}, error)
	Create(object map[string]interface{}) error
	Update(id uint, object map[string]interface{}) error
	Find(id uint) error
	Delete(id uint) error
}
