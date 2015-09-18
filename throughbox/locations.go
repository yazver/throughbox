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

package main

import (
	"strings"
	"path/filepath"
	"github.com/yazver/throughbox/libs/osutils"
	log "github.com/Sirupsen/logrus"
)

const (
	LocationConfigFile = iota
	LocationLogFile 	
)

var Locations map[int] string = {
	LocationConfigFile: "${config}/config.json"
	LocationLogFile: "${config}/throughbox.log"
}

func formatLocation(location, paramName, param string) string {
	return filepath.Join(strings.Split(strings.Replace(location, "${" + paramName + "}", param, -1), "/"))
}

func initLocations() {
	if appConfigDir, err := osutils.GetAppDir(); err != nil {
		log.Warnln("Get application directory: " + err.Error())
	}
	if err != nil || !osutils.PathExists(formatLocation(Location[LocationConfigFile], "config", appDir)) {
		if appConfigDir, err = osutils.GetAppConfigDir("throughbox"); err != nil {
			log.Errorln("Get config directory: " + err.Error())	
		}
	}
	if strings.TrimSpace(appConfigDir) == "" {
		log.Fatalln("The config path is not defined.")
	}
	
	for key, location := range Locations {
		Locations[key] = formatLocation(location, "config", appConfigDir)
	}
}