package config

type Config struct {
	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Zap     Zap     `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
	Mysql   Mysql   `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
}
