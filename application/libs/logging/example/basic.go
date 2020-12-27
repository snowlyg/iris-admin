package example

//var (
//	log1 *logging.Logger //default options, default output is stdout
//	log2 *logging.Logger
//	log3 *logging.Logger //disable color, disable FullTimestamp
//	log4 *logging.Logger //disable level
//	test1 *logging.Logger //disable level
//)
//
//func init() {
//	//default output is stdout
//	//default level is out
//	//default timestamp format is TIMENANO = "2006-01-02 15:04:05.9999999"
//	//default enable colors
//	//default print level and function info
//	//if you need details, see options.go file
//	log1 = logging.NewLogger(&logging.Options{
//	})
//
//	//using hourly rolling,	logging support DAILY and HOURLY rolling
//	//using TIMESECOND = "2006-01-02 15:04:05" timestamp format
//	//set level to info
//	//there are five logging level
//	//"debug","info","warning","error","fatal", debug is default value
//	log2 = logging.NewLogger(&logging.Options{
//		Rolling: logging.DAILY,
//		TimesFormat: logging.TIMESECOND,//default is TIMENANO, you can also use time.RFC3339,eg...
//		Level: "info",
//	})
//
//	//disable colors, and print time passed since beginning of execution
//	log3 = logging.NewLogger(&logging.Options{
//		DisableColors: true,
//		DisableFullTimestamp: true,
//	})
//
//	//do not print level and function info
//	log4 = logging.NewLogger(&logging.Options{
//		DisableLevel: true,
//	})
//
//	// create more than one logger at one time
//	// those logger share the same options
//	// you can use Log("test1"),Log("test3"),Log("test3") to require Logger
//	test1 = logging.NewLogger(&logging.Options{
//	}, "test1.log", "test2.log", "test3.log")
//}
//
//func main() {
//	log1.Debugf("This is default configuration logger...")
//
//	log2.Infof("TIMESECOND, hourly roiing logger...")
//	log3.Infof("DisableColors DisableFullTimestamp logger...")
//	time.Sleep(time.Second)
//	log3.Infof("DisableColors DisableFullTimestamp logger...")
//	log4.Infof("DisableLevel logger...")
//
//	//require test1.log Logger
//	test1.Infof("logging to test1.log")
//
//	//require test3.log Logger
//	logging.Log("test3").Infof("logging to test3.log")
//}
//
