package validate

import (
	"encoding/json"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/viper_server"
	"github.com/spf13/viper"
)

var CONFIG = Config{
	Locale: "zh",
}

type Config struct {
	Locale string `mapstructure:"locale" json:"locale" yaml:"locale"`
}

// IsExist config file is exist
func IsExist() bool {
	return getViperConfig().IsFileExist()
}

// Remove remove config file
func Remove() error {
	return getViperConfig().RemoveFile()
}

// Recover
func Recover() error {
	b, err := json.MarshalIndent(CONFIG, "", "\t")
	if err != nil {
		return err
	}
	return getViperConfig().Recover(b)
}

// getViperConfig get viper config
func getViperConfig() viper_server.ViperConfig {
	configName := "validate"
	return viper_server.ViperConfig{
		Directory: g.ConfigDir,
		Name:      configName,
		Type:      g.ConfigType,
		Watch: func(vi *viper.Viper) error {
			if err := vi.Unmarshal(&CONFIG); err != nil {
				return fmt.Errorf("get Unarshal error: %v", err)
			}
			// watch config file change
			vi.OnConfigChange(func(e fsnotify.Event) {
				fmt.Println("Config file changed:", e.Name)
			})
			vi.WatchConfig()
			return nil
		},
		//
		Default: []byte(`
{
	"locale": "` + CONFIG.Locale + `"
}`),
	}
}
