package e

import (
	"errors"
)

var (
	ErrEmptyName        = errors.New("config'name can't be empty value")
	ErrViperConfInvalid = errors.New("viper conf not invalid")
	ErrConfigInvalid    = errors.New("config not invalid")
)
