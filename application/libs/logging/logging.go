package logging

import (
	"path"

	"github.com/snowlyg/iris-admin/application/libs"
)

var DebugLogger *Logger
var ErrorLogger *Logger
var InfoLogger *Logger

func init() {
	DebugLogger = NewLogger(&Options{
		Level:       "debug",
		Rolling:     DAILY,
		TimesFormat: TIMESECOND,
	}, path.Join(libs.LogDir(), "./debug.log"))
	DebugLogger.SetLogPrefix("log_prefix")

	ErrorLogger = NewLogger(&Options{
		Level:       "error",
		Rolling:     DAILY,
		TimesFormat: TIMESECOND,
	}, path.Join(libs.LogDir(), "./error.log"))
	ErrorLogger.SetLogPrefix("log_prefix")

	InfoLogger = NewLogger(&Options{
		Level:       "info",
		Rolling:     DAILY,
		TimesFormat: TIMESECOND,
	}, path.Join(libs.LogDir(), "./info.log"))
	InfoLogger.SetLogPrefix("log_prefix")
}
