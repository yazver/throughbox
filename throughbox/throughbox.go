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
	//"fmt"
	"flag"
	"os"

	log "github.com/Sirupsen/logrus"
)

var throughBox *ThroughBox = NewThroughBox()

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {

	var (
		configDir   string
		hideConsole bool
		showConsole bool
		debugLog    bool
	)

	flag.StringVar(&configDir, "configdir", "", "Configuration directory")
	flag.BoolVar(&hideConsole, "hidecon", false, "Hide console")
	flag.BoolVar(&showConsole, "showcon", false, "Show console")
	flag.BoolVar(&debugLog, "debuglog", false, "Show console")

	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)

    if configDir != "" {
        Locations[LocationConfigFile] = configDir
    }        
	InitLocations()
	throughBox.LoadConfig(Locations[LocationConfigFile], true)

	//	for _, item := range PortMapList {
	//		fmt.Printf("%#v %#v %#v \n", item.Port, *(item.SourceIP), item.DestinationAdress)
	//	}

	log.Debugln("Start")
	throughBox.Start()
	throughBox.Wait()
}
