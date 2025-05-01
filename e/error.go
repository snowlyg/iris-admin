package e

import (
	"errors"
)

var (
	ErrConfigNameEmpty  = errors.New("config file name empty")
	ErrViperConfInvalid = errors.New("viper conf not invalid")
	ErrConfigInvalid    = errors.New("config not invalid")

	ErrAuthInvalid      = errors.New("auth invalid")
	ErrDbTableNameEmpty = errors.New("database table name empty")
)
