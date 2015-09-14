//   Copyright (C) 2015 Evgeny M. Safonov
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

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
