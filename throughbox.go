package main

import (
	"fmt"
	"os"
	"sync"
	log "github.com/Sirupsen/logrus"
	"github.com/go-fsnotify/fsnotify"
	"./core"

)

var ConfigPath = "./config.json"
var PortMapList core.PortMapList

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {
	log.SetOutput(os.Stderr)

	var err error
	PortMapList, err = core.NewPortMapListFromFile(ConfigPath)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	for _, item := range PortMapList {
		fmt.Printf("%#v %#v %#v \n", item.Port, *(item.SourceIP), item.DestinationAdress)
	}

	log.Debug("Start")
	wait := &sync.WaitGroup{}
	PortMapList.Start(wait)
	wait.Wait()
}
