package seed

type SeedFunc interface {
	Init() (err error)
}

// Seed 数据填充
func Seed(SeedFunctions ...SeedFunc) error {
	if len(SeedFunctions) == 0 {
		return nil
	}
	for _, v := range SeedFunctions {
		err := v.Init()
		if err != nil {
			return err
		}
	}
	return nil
}
