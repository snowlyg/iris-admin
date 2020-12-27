package logging

import (
	_ "log"
	_ "net/http/pprof"
	_ "os/signal"
)

//
//func TestOldLogger(t *testing.T) {
//
//	l := New()
//	l.SetRotateByHour()
//	//	l.SetPrintLevel(false)
//	l.SetHighlighting(false)
//	l.Infof("hahahah")
//	time.Sleep(time.Second * 2)
//	l.Debugf("hasdhdhash %s", 123)
//
//	log := l.Logger()
//	log.Printf("this logger")
//	time.Sleep(time.Second * 2)
//	l.SetLogPrefix(fmt.Sprintf("%d|", os.Getpid()))
//	l.SetOutputByName("./testold.log")
//	l.Info("this logger 2222")
//	l.SetFlags(0)
//	l.Info("this logger 2222")
//	l.Info("this logger 2222")
//	//	os.Remove("./testold.log")
//}
//
//func TestLoggerWithFiled(t *testing.T) {
//	l := With("test", "test value", "key", "1", "value", "1")
//	l.Debugw("hahhh logw url", "test", 1234)
//	l.Debugf("hahhh logf %s url %d", "test", 1234)
//
//	Debugw("debugw test message", "url", "http://service.inke.cn/serviceinfo", "timeout", 3, "retry", 10)
//
//}
//
//func TestDataLogger(t *testing.T) {
//	InitData("./bigdata/trans.log", DailyRolling)
//	DataLog("topic_test", "url", "http://service.inke.cn/serviceinfo/info", "timeout", 3, "retry", 10)
//	DataLog("topic_test", "url", "http://service.inke.cn/serviceinfo/info", "timeout", 3, "retry", 10, "info", map[string]interface{}{"key": "value", "key2": "value2"})
//}
//
//func TestDataLoggerWithKey(t *testing.T) {
//	InitDataWithKey("./bigdata/trans.log", DailyRolling, "test_bigdata")
//	DataLog("topic_test", "url", "http://service.inke.cn/serviceinfo/info", "timeout", 3, "retry", 10)
//	DataLog("topic_test", "url", "http://service.inke.cn/serviceinfo/info", "timeout", 3, "retry", 10, "info", map[string]interface{}{"key": "value", "key2": "value2"})
//}
//
//func TestNewLogger(t *testing.T) {
//	l := NewLogger(&Options{
//		//		DisableColors: true,
//		//		DisableLevel: true,
//		Rolling: SECONDLY,
//		//		TimesFormat: time.RFC3339Nano,
//		TimesFormat: TIMESECOND,
//	}, "test1.log", "test2.log")
//
//	l.SetLogPrefix("log_prefix")
//	//l.SetOutputByName("./test.log")
//	//	l.SetRotateBySecond()
//	//	l.Info("hahahah")
//	l.Debugf("hasdhdhash %d", 123)
//
//	for i := 0; i < 700; i++ {
//		Log("test2").Infof("hahahah %d", i)
//	}
//	//os.Remove("./test.log*")
//}
//
//func BenchmarkDebugLogParallelZap(b *testing.B) {
//	fileobj, _ := os.OpenFile("test3.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
//
//	ec := zap.NewProductionEncoderConfig()
//	ec.EncodeDuration = zapcore.NanosDurationEncoder
//	ec.EncodeTime = zapcore.EpochNanosTimeEncoder
//	enc := zapcore.NewJSONEncoder(ec)
//	sugar := zap.New(zapcore.NewCore(
//		enc,
//		fileobj,
//		zap.DebugLevel,
//	))
//
//	b.ResetTimer()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			//Debugf("test log debug %d", 10234)
//			sugar.Info("test log debug")
//			//fileobj.WriteString(fmt.Sprintf("test log debug %d", 1234))
//		}
//	})
//}
//
//func BenchmarkDebugLogParallel(b *testing.B) {
//	//fileobj, _ := os.OpenFile("test3.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
//	SetOutputByName("test3.log")
//	b.ResetTimer()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			//Debugf("test log debug %d", 10234)
//			Info("test log debug")
//			//fileobj.WriteString(fmt.Sprintf("test log debug %d", 1234))
//		}
//	})
//}
//
//func BenchmarkDebugLogParallelLogrus(b *testing.B) {
//	fileobj, _ := os.OpenFile("test4.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
//	logrusLog := logrus.New()
//	logrusLog.Out = fileobj
//	b.ResetTimer()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			//Debugf("test log debug %d", 10234)
//			//fileobj.WriteString(fmt.Sprintf("test log debug %d", 1234))
//			//	logrusLog.Infof("test log debug")
//			logrusLog.WithFields(logrus.Fields{"adsaf": "fsfsd"}).Info("test log debug")
//		}
//	})
//}
//
//func BenchmarkDebugLog1(b *testing.B) {
//	logrusLog := logrus.New()
//	file, err := os.OpenFile("test.log", os.O_RDWR, 0660)
//	if err != nil {
//		b.Fatal(err)
//	}
//	logrusLog.Out = file
//	for i := 0; i < b.N; i++ {
//		logrusLog.Info("test log debug")
//	}
//}
//
//func BenchmarkDebugLog(b *testing.B) {
//	SetOutputByName("test.log")
//	for i := 0; i < b.N; i++ {
//		//Debugf("test log debug")
//		Infof("test log debug")
//	}
//}
//
//func TestDefaultLog(t *testing.T) {
//	//InitError("log")
//	log := NewLogger(&Options{})
//	log.Errorf("%s", "this is error")
//}
//
//func setUpLogger() {
//	cc := CommonLogConfig{
//		Pathprefix:      "logs/",
//		Rotate:          "day",
//		GenLogLevel:     "info",
//		BalanceLogLevel: "info",
//	}
//	SetOutputPath("logs/")
//	InitCommonLog(cc)
//
//}
//func BenchmarkCommonErrorLog(b *testing.B) {
//	setUpLogger()
//	b.ReportAllocs()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			Errorf("test log debug")
//		}
//	})
//}
//func BenchmarkCommonDebugLog(b *testing.B) {
//	setUpLogger()
//	b.ReportAllocs()
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			Debugf("test log debug")
//		}
//	})
//}

// func BenchmarkDebugLog(b *testing.B) {
// 	SetOutputByName("test.log")
// 	for i := 0; i < b.N; i++ {
// 		Debugf("test log debug")
// 	}
// }

// func BenchmarkDebugLogParallel(b *testing.B) {
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			Debugf("test log debug %d", 10234)
// 		}
// 	})
// }
