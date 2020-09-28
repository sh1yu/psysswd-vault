package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	defaultVaultConfigFile  = "./config.yaml"
	defaultPersistMetaFile  = "meta.data"
	defaultPersistIndexFile = "index.data"
	defaultPersistDataFile  = "file.data"
)

type VaultConfig struct {
	UserConf    UserConfig    `yaml:"user"`
	PersistConf PersistConfig `yaml:"persist"`
}

type UserConfig struct {
	DefaultUserName string `yaml:"defaultUserName"`
}

type PersistConfig struct {
	MetaFile  string `yaml:"meta_path"`
	IndexFile string `yaml:"index_path"`
	DataFile  string `yaml:"data_path"`
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

	if vaultConf.PersistConf.MetaFile == "" {
		vaultConf.PersistConf.MetaFile = defaultPersistMetaFile
	}
	if vaultConf.PersistConf.IndexFile == "" {
		vaultConf.PersistConf.IndexFile = defaultPersistIndexFile
	}
	if vaultConf.PersistConf.DataFile == "" {
		vaultConf.PersistConf.DataFile = defaultPersistDataFile
	}
	return &vaultConf, nil
}
