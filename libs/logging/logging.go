package logging

var Dbug *Logger
var Err *Logger

func init() {
	Dbug = NewLogger(&Options{
		Rolling:     DAILY,
		TimesFormat: TIMESECOND,
	}, "./logs/debug.log")
	Dbug.SetLogPrefix("log_prefix")

	Err = NewLogger(&Options{
		Rolling:     DAILY,
		TimesFormat: TIMESECOND,
	}, "./logs/error.log")
	Err.SetLogPrefix("log_prefix")
}
