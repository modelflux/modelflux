package model

import "github.com/spf13/viper"

type Model interface {
	New(cfg *viper.Viper) error
	Generate(input string) (string, error)
}
