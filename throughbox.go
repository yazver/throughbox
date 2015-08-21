package main

import (
	"fmt"
	"os"
	"sync"
	"log"
	"github.com/go-fsnotify/fsnotify"
	"./core"
	"./debug"
)

var ConfigPath = "./settings.txt"
var PortMapList core.PortMapList

func WatchConfigChanges() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	//defer watcher.Close()

	//done := make(chan bool)
	go func() {
		defer watcher.Close()
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Modified config file: ", event.Name)
					PortMapList.Load(ConfigPath)
					//watcher.Add(ConfigPath)
					for _, item := range PortMapList {
						fmt.Printf("%#v %#v %#v \n", item.Port, *(item.SourceIP), item.DestinationAdress)
					}
				}
			case err := <-watcher.Errors:
				log.Println("Watch config error:", err)
			}
		}
		debug.Println("Stop watch")
	}()

	err = watcher.Add(ConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	//<-done
}

func main() {
	debug.SetDebugMode(debug.DConsole)
	log.SetOutput(os.Stderr)
	
	var err error
	PortMapList, err = core.NewPortMapListFromFile(ConfigPath)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	WatchConfigChanges()

	for _, item := range PortMapList {
		fmt.Printf("%#v %#v %#v \n", item.Port, *(item.SourceIP), item.DestinationAdress)
	}

	debug.Print("Start")
	wait := &sync.WaitGroup{}
	PortMapList.Start(wait)
	wait.Wait()
}
