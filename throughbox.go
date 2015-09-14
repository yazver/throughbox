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
	"os"
	log "github.com/Sirupsen/logrus"
	"github.com/yazver/throughbox/core"
)

const ConfigPath = "./config.json"
var throughBox *core.ThroughBox = core.NewThroughBox()

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)

	throughBox.LoadConfig(ConfigPath, true)

//	for _, item := range PortMapList {
//		fmt.Printf("%#v %#v %#v \n", item.Port, *(item.SourceIP), item.DestinationAdress)
//	}

	log.Debugln("Start")
	throughBox.Start()
	throughBox.Wait()
}
