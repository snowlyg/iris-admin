package seed

type SeedFunc interface {
	Init() (err error)
}

// Seed 数据填充
func Seed(SeedFunctions ...SeedFunc) error {
	for _, v := range SeedFunctions {
		err := v.Init()
		if err != nil {
			return err
		}
	}
	return nil
}
