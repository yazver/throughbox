package core

import (
	"os"
	"encoding/json"
	"errors"
	"fmt"
)

type ConfigPortMap struct {
	SourceIP    string `json:sourceIP` // 192.168.0.100/24
	Port        uint    `json:port`     // 9000
	Destination string `json:dest`     // bbc.com:80
	ACL         string `json:acl`
}
type Config struct {
	SaveLog          bool `json:saveLog`
	SaveErrorLog     bool `json:saveErrorLog`
	ShowLogInConsole bool `json:showLogInConsole`
	PortMap          [] ConfigPortMap `json:portmap`
}

func LoadConfig(filePath string) (*Config, error) {

	config := &Config{}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New(fmt.Errorf("Don't possible open config file (%s): (%s) \n", filePath, err.Error()))
	}
	if err := json.Unmarshal(json.NewDecoder(file), config); err != nil {
		return nil, errors.New(fmt.Errorf("Don't possible decode config file (%s): (%s) \n", filePath, err.Error()))
	}

	return config, nil

}