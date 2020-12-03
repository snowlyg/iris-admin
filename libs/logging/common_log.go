package logging

import (
	"strings"
)

var slowlog *Logger = New()
var genlog *Logger = New()
var crashlog *Logger = New()
var balancelog *Logger = New()

var (
	DAY_ROTATE  = "day"
	HOUR_ROTATE = "hour"
)

type CommonLogConfig struct {
	Pathprefix      string
	Rotate          string
	GenLogLevel     string
	BalanceLogLevel string
}

var isInit bool = false
var trunon bool = true

func init() {

	isInit = false
}

func isHourRotate(rotate string) bool {

	if rotate == HOUR_ROTATE {
		return true
	}
	return false
}

func CloseCommonLog() {
	trunon = false
}

func OpenCommonLog() {
	trunon = true
}

func getNewPathName(clc CommonLogConfig, name string) string {

	hasSuf := strings.HasSuffix(clc.Pathprefix, "/")
	if hasSuf {
	} else {
		name = "/" + name
	}
	return clc.Pathprefix + name
}

func setCommonRotate(clc CommonLogConfig) {
	isHour := isHourRotate(clc.Rotate)
	if isHour {
		slowlog.SetRotateByHour()
		genlog.SetRotateByHour()
		balancelog.SetRotateByHour()
		crashlog.SetRotateByHour()
	} else {
		slowlog.SetRotateByDay()
		genlog.SetRotateByDay()
		balancelog.SetRotateByDay()
		crashlog.SetRotateByDay()
	}
}

func setCommonLogLevel(clc CommonLogConfig) {

	slowlog.SetLevelByString("debug")
	genlog.SetLevelByString(clc.GenLogLevel)
	crashlog.SetLevelByString("debug")
	balancelog.SetLevelByString(clc.BalanceLogLevel)
}

func setCommonOutput(clc CommonLogConfig) {
	slowlog.SetOutputByName(getNewPathName(clc, "slow"))
	genlog.SetOutputByName(getNewPathName(clc, "gen"))
	crashlog.SetOutputByName(getNewPathName(clc, "crash"))
	balancelog.SetOutputByName(getNewPathName(clc, "balance"))
}

func setCommonTimer(clc CommonLogConfig) {

	balancelog.SetPrintLevel(false)
	slowlog.SetPrintLevel(false)
	genlog.SetPrintLevel(false)
	crashlog.SetPrintLevel(false)

	balancelog.SetTimeFmt(TIMEMICRO)
	slowlog.SetTimeFmt(TIMEMICRO)
	genlog.SetTimeFmt(TIMEMICRO)
	crashlog.SetTimeFmt(TIMEMICRO)

}

func InitCommonLog(clc CommonLogConfig) string {

	if len(clc.Pathprefix) == 0 {
		return ""
	}
	setCommonRotate(clc)
	setCommonOutput(clc)
	setCommonLogLevel(clc)
	setCommonTimer(clc)
	isInit = true
	return ""
}

func checkOpenStatus() bool {

	//kai guan
	if trunon == false {
		return false
	}
	return true
}

func checkNeedLog() bool {

	if isInit == false {
		return false
	}

	return true
}

func SlowLog(v ...interface{}) {

	if !checkOpenStatus() {
		return
	}

	if !checkNeedLog() {
		_defaultLogger.Debug(v...)
		return
	}
	slowlog.Debug(v...)
}

func SlowLogf(format string, v ...interface{}) {

	if !checkOpenStatus() {
		return
	}

	if !checkNeedLog() {
		_defaultLogger.Debugf(format, v...)
		return
	}
	slowlog.Debugf(format, v...)
}

func GenLog(v ...interface{}) {

	if !checkOpenStatus() {
		return
	}

	if !checkNeedLog() {
		_defaultLogger.Debug(v...)
		return
	}

	genlog.Debug(v...)
}

func GenLogf(format string, v ...interface{}) {

	if !checkOpenStatus() {
		return
	}

	if !checkNeedLog() {
		_defaultLogger.Debugf(format, v...)
		return
	}

	genlog.Debugf(format, v...)
}

func CrashLog(v ...interface{}) {

	if !checkOpenStatus() {
		return
	}

	if !checkNeedLog() {
		_defaultLogger.Debug(v...)
		return
	}

	crashlog.Debug(v...)
}

func CrashLogf(format string, v ...interface{}) {

	if !checkOpenStatus() {
		return
	}

	if !checkNeedLog() {
		_defaultLogger.Debugf(format, v...)
		return
	}
	crashlog.Debugf(format, v...)
}

func BalanceLog(v ...interface{}) {

	if !checkOpenStatus() {
		return
	}
	if !checkNeedLog() {
		_defaultLogger.Debug(v...)
		return
	}
	balancelog.Debug(v...)
}

func BalanceLogf(format string, v ...interface{}) {

	if !checkOpenStatus() {
		return
	}
	if !checkNeedLog() {
		_defaultLogger.Debugf(format, v...)
		return
	}
	balancelog.Debugf(format, v...)
}
