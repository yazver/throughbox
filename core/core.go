package core

import (
	"sync"
	"encoding/json"
	"log"
	"os"
)

var ThroughBox *ThroughBox = newThroughBox()

type ThroughBox struct {
	wg            *sync.WaitGroup
	locker        *sync.RWMutex
	pmList        PortMapList

	configPath    string
	configWatcher *ConfigWatcher
}

func newThroughBox() *ThroughBox {
	return &ThroughBox{&sync.WaitGroup{}, &sync.Mutex{}, PortMapList{}}
}

func (tb *ThroughBox) Start() {
	tb.locker.RLock()
	defer tb.locker.RUnlock()

	tb.pmList.Start(tb.wg)
}

func (tb *ThroughBox) Stop() {
	tb.locker.RLock()
	defer tb.locker.RUnlock()

	tb.pmList.Stop()
}

func (tb *ThroughBox) loadConfig() {
	tb.locker.Lock()
	defer tb.locker.Unlock()
	tb.wg.Add(1)
	defer tb.wg.Done()

	if config, err := LoadConfig(tb.configPath); err == nil {
		pmList := NewPortMapList()
		if err := pmList.InitFromConfig(config); err != nil {
			log.Fatalln(err.Error())
		}
		tb.pmList.Stop()
		tb.pmList = pmList
	}

}

func (tb *ThroughBox) LoadConfig(configPath string, watchChanges bool) {
	tb.configPath = configPath
	tb.loadConfig()
	if watchChanges {
		if tb.configWatcher != nil {
			tb.configWatcher.Close()
			tb.configWatcher = nil
		}
		tb.configWatcher = NewConfigWatcher(configPath, tb.loadConfig())
	}
}



