package mongodb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/viper_server"
	"github.com/spf13/viper"
)

// init  initialize
func init() {
	viper_server.Init(getViperConfig())
}

var CONFIG = MongoDB{
	DB:      "mongo_test",
	Timeout: 10,
	Addr:    "localhost:27017",
}

type MongoDB struct {
	Timeout time.Duration `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	DB      string        `mapstructure:"db" json:"db" yaml:"db"`
	Addr    string        `mapstructure:"addr" json:"addr" yaml:"addr"`
}

func (md *MongoDB) GetApplyURI() string {
	return str.Join("mongodb://", md.Addr)
}

// IsExist config file is exist
func IsExist() bool {
	return getViperConfig().IsFileExist()
}

// Remove remove config file
func Remove() error {
	return getViperConfig().Remove()
}

// Recover
func Recover() error {
	b, err := json.Marshal(CONFIG)
	if err != nil {
		return err
	}
	return getViperConfig().Recover(b)
}

// getViperConfig get viper config
func getViperConfig() viper_server.ViperConfig {
	configName := "mongo"
	return viper_server.ViperConfig{
		Debug:     true,
		Directory: g.ConfigDir,
		Name:      configName,
		Type:      g.ConfigType,
		Watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(&CONFIG); err != nil {
				return fmt.Errorf("get Unarshal error: %v", err)
			}
			// watch config file change
			vi.SetConfigName(configName)
			return nil
		},
		//
		Default: []byte(`
{
	"timeout": "` + CONFIG.Timeout.String() + `",
	"db": "` + CONFIG.DB + `",
	"addr": "` + CONFIG.Addr + `"
}`),
	}
}
