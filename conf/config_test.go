package conf

import "testing"

func TestSetDefaultAddrAndTimeFormat(t *testing.T) {
	dc := &Conf{}
	addr := ""
	if dc.System.Addr != addr {
		t.Errorf("config system addr want '%s' but get '%s'", addr, dc.System.Addr)
	}
	timeFormat := ""
	if dc.System.TimeFormat != timeFormat {
		t.Errorf("config system time format want '%s' but get '%s'", timeFormat, dc.System.TimeFormat)
	}
	dc.SetDefaultAddrAndTimeFormat()
	addr = "127.0.0.1:80"
	if dc.System.Addr != addr {
		t.Errorf("config system addr want '%s' but get '%s'", addr, dc.System.Addr)
	}
	timeFormat = "2006-01-02 15:04:05"
	if dc.System.TimeFormat != timeFormat {
		t.Errorf("config system time format want '%s' but get '%s'", timeFormat, dc.System.TimeFormat)
	}

	c := NewConf()
	if c.IsExist() {
		t.Error("config exist before init")
	}
	addr = "127.0.0.1:80"
	if c.System.Addr != addr {
		t.Errorf("config system addr want '%s' but get '%s'", addr, c.System.Addr)
	}
	timeFormat = "2006-01-02 15:04:05"
	if c.System.TimeFormat != timeFormat {
		t.Errorf("config system time format want '%s' but get '%s'", timeFormat, c.System.TimeFormat)
	}

	if err := c.Recover(); err != nil {
		t.Error(err.Error())
	}
	defer func() {
		if err := c.getViperConfig().RemoveDir(); err != nil {
			t.Error(err.Error())
		}
		c.RemoveRbacModel()
	}()
	if !c.IsExist() {
		t.Error("config not exist after recover")
	}
	c.RemoveFile()
	if c.IsExist() {
		t.Error("config exist after remove")
	}
	wantUri := "mongodb://localhost:27017"
	if c.Mongo.GetApplyURI() != wantUri {
		t.Errorf("config mongodb uri want '%s' but get '%s'", wantUri, c.Mongo.GetApplyURI())
	}
}
