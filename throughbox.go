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
