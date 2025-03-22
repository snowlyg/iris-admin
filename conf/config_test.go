package conf

// func TestIsExist(t *testing.T) {
// 	t.Run("test web config IsExist function", func(t *testing.T) {
// 		if !IsExist() {
// 			t.Errorf("config's files is not exist.")
// 		}
// 	})

// 	t.Run("Test Remove function", func(t *testing.T) {
// 		if err := Remove(); err != nil {
// 			t.Error(err)
// 		}
// 		if IsExist() {
// 			t.Errorf("config's files remove is fail.")
// 		}
// 	})
// }

// func TestSetDefaultAddrAndTimeFormat(t *testing.T) {
// 	CONFIG.System.Addr = ""
// 	CONFIG.System.TimeFormat = ""
// 	t.Run("test set defualt addr and time format", func(t *testing.T) {
// 		setDefaultAddrAndTimeFormat()
// 		if CONFIG.System.Addr != "127.0.0.1:80" {
// 			t.Errorf("applyURI want %s but get %s", "127.0.0.1:80", CONFIG.System.Addr)
// 		}
// 		if CONFIG.System.TimeFormat != "2006-01-02 15:04:05" {
// 			t.Errorf("applyURI want %s but get %s", "2006-01-02 15:04:05", CONFIG.System.TimeFormat)
// 		}
// 	})
// }

// func TestToStaticUrl(t *testing.T) {
// 	setDefaultAddrAndTimeFormat()
// 	CONFIG.System.StaticPrefix = "/admin"
// 	t.Run("test to static url with tls", func(t *testing.T) {
// 		CONFIG.System.Tls = true
// 		staticPath := toStaticUrl("/uploads/123.png")
// 		if staticPath != "https://127.0.0.1:80/admin/uploads/123.png" {
// 			t.Errorf("applyURI want %s but get %s", "https://127.0.0.1:80/admin/uploads/123.png", staticPath)
// 		}
// 	})
// 	t.Run("test to static url with tls", func(t *testing.T) {
// 		CONFIG.System.Tls = false
// 		staticPath := toStaticUrl("/uploads/123.png")
// 		if staticPath != "http://127.0.0.1:80/admin/uploads/123.png" {
// 			t.Errorf("applyURI want %s but get %s", "https://127.0.0.1:80/admin/uploads/123.png", staticPath)
// 		}
// 	})
// }
