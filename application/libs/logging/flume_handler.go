package logging

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type FlumeLog struct {
	clientIp string
	topicMap map[string][]string

	loggerMap map[string]*Logger
}

func (self *FlumeLog) readList(configPath string) (map[string][]string, error) {

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf := bufio.NewReader(file)

	var data map[string][]string
	var section string

	data = make(map[string][]string)

	for {

		l, err := buf.ReadString('\n')
		line := strings.TrimSpace(l)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}

			if len(line) == 0 {
				break
			}
		}

		switch {

		case len(line) == 0:
		case string(line[0]) == "#": //增加配置文件备注
		case line[0] == '[' && line[len(line)-1] == ']':

			section = strings.TrimSpace(line[1 : len(line)-1])
			var list []string
			data[section] = list

		default:
			i := strings.IndexAny(line, "=")
			value := strings.TrimSpace(line[0:i])
			if list, ok := data[section]; ok {
				list = append(list, value)
				data[section] = list
			}
		}
	}
	return data, nil
}

func (self *FlumeLog) mockMsg() map[string][]string {

	topicMap := make(map[string][]string)

	var list []string
	list = append(list, "key1")
	list = append(list, "key2")

	topicMap["demo_topic"] = list

	return topicMap
}

func (self *FlumeLog) parseConfig(configPath string) error {

	//通过配置文件,构建,topicMap 出来。
	topicMap, err := self.readList(configPath)

	if err != nil {
		return err
	}
	self.topicMap = topicMap

	return nil
}

func (self *FlumeLog) initLogger() error {

	loggerMap := make(map[string]*Logger)

	for topic, _ := range self.topicMap {

		//	dirName := "/Users/zhaolingzhi/xxx/logs/big_data_flume_log/" + topic + "_bussiness"
		dirName := "/a8root/logs/big_data_flume_log/" + topic + "_bussiness"
		mkErr := Mkdirlog(dirName)

		if mkErr != nil {
			//创建目录失败。
			return mkErr
		}

		fileName := dirName + "/" + topic
		logger := New()
		logger.SetHighlighting(false)
		logger.SetRotateByDay()
		logger.SetPrintLevel(false)
		logger.SetFlags(0)
		logger.SetOutputByName(fileName)
		loggerMap[topic] = logger
	}

	self.loggerMap = loggerMap
	return nil
}

func (self *FlumeLog) InitFlumeLog(configPath string) (int, error) {

	//需要init配置文件。

	err := self.parseConfig(configPath)
	if err != nil {
		return -1, err
	}

	self.initLogger()

	clientIp, errIp := GetInternal()
	if errIp != nil {
		clientIp = "127.0.0.1"
		fmt.Println(errIp)

	}
	self.clientIp = clientIp

	return 0, nil
}

func (sefl *FlumeLog) makeBussinessName(topic string) string {
	return topic
}

func (self *FlumeLog) makeBussinessInfo(topic string, bussinessValue map[string]string) (string, string) {

	bussinessInfo := ""

	if valueList, ok := self.topicMap[topic]; ok {
		// 找到了

		addSubSize := 0
		totalSubSize := 14 //总共要添加14个\002
		for _, value := range valueList {

			//通过value,去bussinessValue
			if bussValue, ok := bussinessValue[value]; ok {
				//找到了具体的值。
				bussinessInfo = bussinessInfo + bussValue
			}
			bussinessInfo = bussinessInfo + "\002"
			addSubSize = addSubSize + 1
		}

		lessAddSubSize := totalSubSize - addSubSize

		for i := 0; i < lessAddSubSize; i++ {
			bussinessInfo = bussinessInfo + "\002"
		}

	} else {
		errMsg := "no found topic:" + topic
		return "", errMsg
	}
	return bussinessInfo, ""
}

func (self *FlumeLog) makeAtomicInfo(uid int64, atomic string) string {

	//支持 kv格式 111=222&222=333

	atomicArr := strings.Split(atomic, "&")

	atomicMap := make(map[string]string)

	for _, atomicValue := range atomicArr {
		valueArr := strings.Split(atomicValue, "=")
		valueLen := len(valueArr)
		if valueLen != 2 {
			continue
		}

		key := valueArr[0]
		value := valueArr[1]
		atomicMap[key] = value
	}

	atomicInfo := ""
	//开始拼接字符串。
	//var lc, cc, cv, ua, devi, imsi, imei, osversion, conn, proto, tg, smid, client_ip, idfa, aid, appid ,logid , mjid string
	//var ok bool

	uidStr := strconv.FormatInt(uid, 10)
	atomicInfo = atomicInfo + uidStr + "\002"
	if lc, ok := atomicMap["lc"]; ok {
		//存在
		atomicInfo = atomicInfo + lc + "\002"
	} else {
		atomicInfo = atomicInfo + "\002"
	}

	if cc, ok := atomicMap["cc"]; ok {
		//存在
		atomicInfo = atomicInfo + cc + "\002"
	} else {
		atomicInfo = atomicInfo + "\002"
	}

	if cv, ok := atomicMap["cv"]; ok {
		//存在
		atomicInfo = atomicInfo + cv + "\002"
	} else {
		atomicInfo = atomicInfo + "\002"
	}
	if ua, ok := atomicMap["ua"]; ok {
		//存在
		atomicInfo = atomicInfo + ua + "\002"
	} else {
		atomicInfo = atomicInfo + "\002"
	}
	if devi, ok := atomicMap["devi"]; ok {
		//存在
		atomicInfo = atomicInfo + devi + "\002"
	} else {
		atomicInfo = atomicInfo + "\002"
	}
	if imsi, ok := atomicMap["imsi"]; ok {
		//存在
		atomicInfo = atomicInfo + imsi + "\002"
	} else {
		atomicInfo = atomicInfo + "\002"
	}
	if imei, ok := atomicMap["imei"]; ok {
		//存在
		atomicInfo = atomicInfo + imei + "\002"
	} else {
		atomicInfo = atomicInfo + "\002"
	}
	if osversion, ok := atomicMap["osversion"]; ok {
		//存在
		atomicInfo = atomicInfo + osversion + "\002"
	} else {
		atomicInfo = atomicInfo + "\002"
	}
	if conn, ok := atomicMap["conn"]; ok {
		//存在
		atomicInfo = atomicInfo + conn + "\002"
	} else {
		atomicInfo = atomicInfo + "\002"
	}
	if proto, ok := atomicMap["proto"]; ok {
		//存在
		atomicInfo = atomicInfo + proto + "\002"
	} else {
		atomicInfo = atomicInfo + "\002"
	}
	if tg, ok := atomicMap["tg"]; ok {
		//存在
		atomicInfo = atomicInfo + tg + "\002"
	} else {
		atomicInfo = atomicInfo + "\002"
	}
	if smid, ok := atomicMap["smid"]; ok {
		//存在
		atomicInfo = atomicInfo + smid + "\002"
	} else {
		atomicInfo = atomicInfo + "\002"
	}
	if client_ip, ok := atomicMap["client_ip"]; ok {
		//存在
		atomicInfo = atomicInfo + client_ip + "\002"
	} else {
		atomicInfo = atomicInfo + "\002"
	}
	if idfa, ok := atomicMap["idfa"]; ok {
		//存在
		atomicInfo = atomicInfo + idfa + "\002"
	} else {
		atomicInfo = atomicInfo + "\002"
	}
	if aid, ok := atomicMap["aid"]; ok {
		//存在
		atomicInfo = atomicInfo + aid + "\002"
	} else {
		atomicInfo = atomicInfo + "\002"
	}
	if appid, ok := atomicMap["appid"]; ok {
		//存在
		atomicInfo = atomicInfo + appid + "\002"
	} else {
		atomicInfo = atomicInfo + "\002"
	}
	if logid, ok := atomicMap["logid"]; ok {
		//存在
		atomicInfo = atomicInfo + logid + "\002"
	} else {
		atomicInfo = atomicInfo + "\002"
	}
	if mjid, ok := atomicMap["mjid"]; ok {
		//存在
		atomicInfo = atomicInfo + mjid + "\002"
	} else {
		atomicInfo = atomicInfo + "\002"
	}

	atomicInfo = atomicInfo + "\002" + "\002" + "\002" + "\002" + "\002"
	return atomicInfo
}

func (self *FlumeLog) makeTime() string {

	now_time := time.Now().Unix() * 1000
	now_time_string := strconv.FormatInt(now_time, 10)
	recordTimeInfo := now_time_string
	return recordTimeInfo

}
func (self *FlumeLog) makeInterfaceInfo() string {

	interfaceInfo := "\002" + "\002" + "\002" + "\002" + "\002" + "\002" + "\002" + "\002" + "\002"
	return interfaceInfo
}

func (self *FlumeLog) writerFlumeLog(topic string, bigDataMsg string) {

	if logger, ok := self.loggerMap[topic]; ok {
		logger.Info(bigDataMsg)
	}
}

func (self *FlumeLog) WriteBussinessLog(topic string, uid int64, atomic string, bussinessValue map[string]string) (int, string) {

	//这里写具体的业务。
	bussinessName := self.makeBussinessName(topic)
	bussinessInfo, biErrMsg := self.makeBussinessInfo(topic, bussinessValue)
	if len(biErrMsg) != 0 {
		return -1, biErrMsg
	}

	atomicInfo := self.makeAtomicInfo(uid, atomic)
	nowTimeStamp := self.makeTime()
	interfaceInfo := self.makeInterfaceInfo()

	bigDataMsg := bussinessName + "\001" + bussinessInfo + "\001" + atomicInfo + "\001" + nowTimeStamp + "\001" + interfaceInfo + "\001" + self.clientIp

	//	fmt.Println(self.clientIp)

	self.writerFlumeLog(topic, bigDataMsg)
	return 0, ""
}
