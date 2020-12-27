package example

//
//import (
//	"github.com/snowlyg/blog/libs/logging"
//	"time"
//)
//
//func main() {
//	flumeLogger := &logging.FlumeLog{
//	}
//	flumeLogger.InitFlumeLog("config")
//	atomic := "lc=11111&cc=ssss"
//
//	bussinessValue := make(map[string]string)
//	bussinessValue["city"] = "city"
//	bussinessValue["peak"] = "peak"
//
//	for i := 0; i < 20 ; i++ {
//		flumeLogger.WriteBussinessLog("user_join_room", 123455, atomic, bussinessValue)
//		time.Sleep(1 * time.Second)
//	}
//}
