package config

type Limit struct {
	Disable bool    `mapstructure:"disable" json:"disable" yaml:"disable"`
	Limit   float64 `mapstructure:"limit" json:"limit" yaml:"limit"`
	Burst   int     `mapstructure:"burst" json:"burst" yaml:"burst"`
}
