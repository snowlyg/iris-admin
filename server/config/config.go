package config

var CONFIG Config

type Config struct {
	MaxSize int64   `mapstructure:"max-size" json:"burst" yaml:"max-size"`
	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Limit   Limit   `mapstructure:"limit" json:"limit" yaml:"limit"`
	Zap     Zap     `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
	Mysql   Mysql   `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
}
