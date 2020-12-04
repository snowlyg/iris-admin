package libs

import (
	"encoding/json"
	"errors"
	"github.com/snowlyg/blog/libs/logging"
)

// SetCache 设置缓存
func SetCache(key string, obj interface{}) error {
	if Config.Cache {
		client := GetRedisClusterClient()
		defer client.Close()
		data, err := json.Marshal(obj)
		if err != nil {
			logging.Err.Errorf("set cache %s read json err : %+v\n", key, err)
			return err
		}
		if _, err := client.Set(key, data); err != nil {
			logging.Err.Errorf("set cache %s set err : %+v\n", key, err)
			return err
		}
	}

	return nil
}

// GetCache 获取缓存
func GetCache(key string, obj interface{}) error {
	if Config.Cache {
		client := GetRedisClusterClient()
		defer client.Close()
		data, err := client.GetKey(key)
		if err != nil {
			logging.Err.Errorf("get cache %s err : %+v\n", key, err)
			return err
		}
		bdata, ok := data.([]byte)
		if !ok {
			return errors.New("数据类型错误")
		}
		if err := json.Unmarshal(bdata, obj); err != nil {
			logging.Err.Errorf("get cache %s unjson err : %+v\n", key, err)
			return err
		}
	}
	return nil
}
