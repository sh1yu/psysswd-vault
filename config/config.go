package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultVaultConfigFile  = "~/.pvlt/config.yaml"
	defaultPersistMetaFile  = "~/.pvlt/meta.data"
	defaultPersistIndexFile = "~/.pvlt/index.data"
	defaultPersistDataFile  = "~/.pvlt/file.data"
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

	confPath, err := CreateFileIfNeeded(configFile)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(confPath)
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

func CreateFileIfNeeded(file string) (string, error) {

	if strings.HasPrefix(file, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		file = home + strings.TrimPrefix(file, "~")
	}

	absPath, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}

	absDir := filepath.Dir(absPath)

	dInfo, err := os.Stat(absDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(absDir, 0755)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	} else {
		if !dInfo.IsDir() {
			return "", errors.New("dir " + absDir + " is not a dir")
		}
	}

	fInfo, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		_, err = os.Create(absPath)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	} else {
		if !fInfo.Mode().IsRegular() {
			return "", errors.New("file " + absPath + " is not a regular file")
		}
	}

	return absPath, nil
}
