package utils

import (
	"encoding/json"
	"io/ioutil"
)

type VaultConfig struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	Namespace    string `json:"namespace"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	LDAPGroup    string `json:"ldap_group"` // Optional
	SecretEngine string `json:"secretengine"`
}

type Config struct {
	Vaults []VaultConfig `json:"vaults"`
}

func LoadConfig(filename string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
