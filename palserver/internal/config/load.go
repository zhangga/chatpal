package config

import (
	"errors"
	"github.com/spf13/viper"
	"os"
)

var (
	ErrNotFound = errors.New("config not found")
	ErrNotFile  = errors.New("config is not a file")
)

func LoadBootstrapConfig(path string) (*ConfBootstrap, error) {
	viper := viper.New()
	if err := loadConfigByFile(path, viper); err != nil {
		return nil, err
	}
	var conf ConfBootstrap
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func loadConfigByFile(path string, viper *viper.Viper) error {
	if len(path) == 0 {
		return ErrNotFound
	}

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return ErrNotFound
	}
	if info.IsDir() {
		return ErrNotFile
	}
	viper.SetConfigType("yaml")
	viper.SetConfigFile(path)
	return viper.ReadInConfig()
}
