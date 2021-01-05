package logging

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func Mkdirlog(dir string) error {
	_, er := os.Stat(dir)
	b := er == nil || os.IsExist(er)
	if !b {
		if err := os.MkdirAll(dir, 0666); err != nil {
			if os.IsPermission(err) {
				fmt.Println("create dir error:", err.Error())
				return err
			}
		}
	}
	return nil
}

func GetInternal() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if strings.HasPrefix(ipnet.IP.String(), "10.") == true || strings.HasPrefix(ipnet.IP.String(), "127.") == true {
					return ipnet.IP.String(), nil
				}
			}
		}
	}
	return "", err
}
