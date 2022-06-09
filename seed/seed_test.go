package seed

import "testing"

type testSeed struct{}

func (ts *testSeed) Init() error {
	return nil
}

func TestSeed(t *testing.T) {
	t.Run("数据填充", func(t *testing.T) {
		err := Seed(&testSeed{})
		if err != nil {
			t.Error(err)
		}
	})
}
