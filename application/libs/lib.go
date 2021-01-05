package libs

import (
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"gopkg.in/ini.v1"
	"net"
	"os"
	"path/filepath"
	"time"
)

// md5
func MD5(str string) string {
	encoder := md5.New()
	encoder.Write([]byte(str))
	return hex.EncodeToString(encoder.Sum(nil))
}

// 当前目录
func CWD() string {
	path, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(path)
}

func LogDir() string {
	dir := filepath.Join(CWD(), "logs")
	EnsureDir(dir)
	return dir
}

func EnsureDir(dir string) (err error) {
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return
		}
	}
	return
}

func IsPortInUse(port int64) bool {
	if conn, err := net.DialTimeout("tcp", net.JoinHostPort("", fmt.Sprintf("%d", port)), 3*time.Second); err == nil {
		conn.Close()
		return true
	}
	return false
}

func init() {
	gob.Register(map[string]interface{}{})
	ini.PrettyFormat = false
}
