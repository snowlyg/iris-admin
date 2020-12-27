package logging

var (
	levelsMap = map[string]bool{
		"debug":   true,
		"info":    true,
		"warning": true,
		"error":   true,
		"fatal":   true,
	}

	rollingMap = map[string]bool{
		DAILY:    true,
		HOURLY:   true,
		SECONDLY: true,
	}
)

const (
	TIMENANO  = "2006-01-02 15:04:05.9999999"
	TIMEMICRO = "2006-01-02 15:04:05.999"

	TIMESECOND = "2006-01-02 15:04:05"

	DAILY    = "20060102"
	HOURLY   = "2006010215"
	SECONDLY = "200601021505"
)

type Options struct {
	// There are five logging level
	// "debug","info","warning","error","fatal", debug is default value
	Level string

	// Set to true to bypass checking for a TTY before outputting colors.
	//ForceColors bool

	Prefix string

	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	DisableFullTimestamp bool

	// TimesFormat to use for display when a full timestamp is printed
	TimesFormat string

	// Whether printf level string when logging or not
	DisableLevel bool

	// Force disabling colors.
	DisableColors bool

	DisableTimestamp bool

	// "daily", "hourly", default is not rolling
	Rolling string

	// This option will not wrap empty fields in quotes if true
	DisableQuoteEmptyFields bool

	// The fields are sorted by default for a consistent output.
	// Not sorting is default value.
	Sorting bool
}

func (self *Options) init() {
	if !levelsMap[self.Level] {
		self.Level = "debug"
	}

	if !rollingMap[self.Rolling] {
		self.Rolling = ""
	}

	if self.TimesFormat == "" {
		self.TimesFormat = TIMEMICRO
	}
}
