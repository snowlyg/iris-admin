package logging

import (
	"github.com/snowlyg/blog/application/libs"
	"path"
)

var DebugLogger *Logger
var ErrorLogger *Logger
var InfoLogger *Logger

func init() {
	DebugLogger = NewLogger(&Options{
		Rolling:     DAILY,
		TimesFormat: TIMESECOND,
	}, path.Join(libs.LogDir(), "./debug.log"))
	DebugLogger.SetLogPrefix("log_prefix")

	ErrorLogger = NewLogger(&Options{
		Rolling:     DAILY,
		TimesFormat: TIMESECOND,
	}, path.Join(libs.LogDir(), "./error.log"))
	ErrorLogger.SetLogPrefix("log_prefix")

	InfoLogger = NewLogger(&Options{
		Rolling:     DAILY,
		TimesFormat: TIMESECOND,
	}, path.Join(libs.LogDir(), "./info.log"))
	InfoLogger.SetLogPrefix("log_prefix")
}
