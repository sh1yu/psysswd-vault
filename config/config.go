package config

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultVaultConfigFile = "~/.pvlt/config.yaml"
	defaultPersistDataFile = "~/.pvlt/file.data"
)

type VaultConfig struct {
	UserConf    UserConfig         `yaml:"user"`
	PersistConf PersistConfig      `yaml:"persist"`
	RemoteConf  RemoteConfig       `yaml:"remote"`
	Credentials []CredentialConfig `yaml:"credentials"`
}

type UserConfig struct {
	DefaultUserName string `yaml:"defaultUserName"`
}

type PersistConfig struct {
	DataFile string `yaml:"data_path"`
}

type RemoteConfig struct {
	ServerAddr string `yaml:"server_addr"`
}

type CredentialConfig struct {
	User  string `yaml:"user"`
	Token string `yaml:"token"`
}

func InitConf(configFile string, err error) (*VaultConfig, error) {
	if err != nil {
		return nil, err
	}
	if configFile == "" {
		configFile = defaultVaultConfigFile
	}

	confPath, err := createFileIfNeeded(configFile)
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

	if vaultConf.PersistConf.DataFile == "" {
		vaultConf.PersistConf.DataFile = defaultPersistDataFile
	}
	return &vaultConf, nil
}

func InitDBFile(dbFile string) (*gorm.DB, error) {
	f, err := createFileIfNeeded(dbFile)
	if err != nil {
		return nil, err
	}
	return gorm.Open("sqlite3", f)
}

func createFileIfNeeded(file string) (string, error) {

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
