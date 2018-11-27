package system

import (
	"fmt"
	"github.com/pelletier/go-toml"
)

/**
 * 返回单例实例
 * @method New
 */
func configNew() *toml.Tree {
	config, err := toml.LoadFile("./config.toml")

	if err != nil {
		fmt.Println("Error ", err.Error())
	}

	return config
}
