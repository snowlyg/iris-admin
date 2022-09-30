package g

import "testing"

func TestCons(t *testing.T) {
	t.Run("test constant param", func(t *testing.T) {
		if ConfigType != "json" {
			t.Errorf("ConfigType want %s but get %s", "json", ConfigType)
		}
		if ConfigDir != "config" {
			t.Errorf("ConfigDir want %s but get %s", "config", ConfigDir)
		}
		if CasbinFileName != "rbac_model.conf" {
			t.Errorf("CasbinFileName want %s but get %s", "rbac_model.conf", CasbinFileName)
		}
		if StatusUnknown != 0 {
			t.Errorf("StatusUnknown want %d but get %d", 0, StatusUnknown)
		}
		if StatusTrue != 1 {
			t.Errorf("StatusTrue want %d but get %d", 1, StatusTrue)
		}
		if StatusFalse != 2 {
			t.Errorf("StatusFalse want %d but get %d", 2, StatusFalse)
		}
	})
}
