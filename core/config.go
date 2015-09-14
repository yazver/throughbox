package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type ConfigPortMap struct {
	SourceIP    string `json:sourceIP` // 192.168.0.100/24
	Port        int    `json:port`     // 9000
	Destination string `json:dest`     // bbc.com:80
	ACL         string `json:acl`
}
type Config struct {
	SaveLog          bool            `json:saveLog`
	LogDebugInfo     bool            `json:logDebugInfo`
	ShowLogInConsole bool            `json:showLogInConsole`
	PortMap          []ConfigPortMap `json:portmap`
}

func LoadConfig(filePath string) (*Config, error) {

	config := &Config{}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Don't possible open config file (%s): (%s) \n", filePath, err.Error()))
	}
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, errors.New(fmt.Sprintf("Don't possible decode config file (%s): (%s) \n", filePath, err.Error()))
	}

	return config, nil

}
