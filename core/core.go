package core

import (
	"io"
	"os"
	"sync"
	log "github.com/Sirupsen/logrus"
)


type ThroughBox struct {
	wg            *sync.WaitGroup
	locker        *sync.RWMutex
	pmList        PortMapList

	configPath    string
	configWatcher *ConfigWatcher
}

func NewThroughBox() *ThroughBox {
	return &ThroughBox{&sync.WaitGroup{}, &sync.RWMutex{}, PortMapList{}, "", nil}
}

func (tb *ThroughBox) Wait() {
	tb.wg.Wait()
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
		// Init base settings
		if config.LogDebugInfo {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.InfoLevel)
		}
		if config.SaveLog && config.ShowLogInConsole {
			log.SetOutput(io.MultiWriter(os.Stderr))
		} else if config.SaveLog {
			
		} else if config.ShowLogInConsole {
			log.SetOutput(os.Stderr)
		} else {
			log.SetOutput(io.MultiWriter())
		}

		// Init port mapping
		pmList := NewPortMapList()
		if err := pmList.InitFromConfig(config); err != nil {
			log.Errorln("Can't init PortMapList fron config: " + err.Error())
		}
		tb.pmList.Stop()
		tb.pmList = pmList
	} else {
		log.Errorln("Can't load config file: " + err.Error())
	}

	for _, item := range tb.pmList {
		log.Debugf("Port map: %s", item)
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
		var err error
		tb.configWatcher, err = NewConfigWatcher(configPath, tb.loadConfig)
		if err != nil {
			log.Errorln(err.Error())
		}
	}
}



