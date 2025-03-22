package admin

import "testing"

func TestSeed(t *testing.T) {
	if err := Seed(&testSeed{}); err != nil {
		t.Errorf("test seed fail:%s", err.Error())
	}
}
