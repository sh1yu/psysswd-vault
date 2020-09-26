package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const defaultVaultConfigFile = "./config.yaml"

type VaultConfig struct {
	UserConf UserConfig `yaml:"user"`
}

type UserConfig struct {
	DefaultUserName string `yaml:"defaultUserName"`
}

func InitConf(configFile string, err error) (*VaultConfig, error) {
	if err != nil {
		return nil, err
	}
	if configFile == "" {
		configFile = defaultVaultConfigFile
	}
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var vaultConf = VaultConfig{}
	err = yaml.Unmarshal(content, &vaultConf)
	if err != nil {
		return nil, err
	}
	return &vaultConf, nil
}
