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
	log "github.com/Sirupsen/logrus"
	"github.com/go-fsnotify/fsnotify"
)

type OnChange func()

type ConfigWatcher struct {
	//	wg *sync.WaitGroup
	//	locker *sync.Mutex
	//	pmList PortMapList
	watcher        *fsnotify.Watcher
	onConfigChange OnChange
	filePath       string
	close          chan struct {}
}

func (configWatcher *ConfigWatcher) Close() {
	close(configWatcher.close)
}

func (configWatcher *ConfigWatcher) Watch() {

	go func() {
		defer configWatcher.watcher.Close()
		for {
			select {
			case <-configWatcher.close:
				return
			case event := <-configWatcher.watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Infoln("Modified config file: ", event.Name)
					configWatcher.onConfigChange()
				}
			case err := <-configWatcher.watcher.Errors:
				log.Errorln("Watch config error:", err)
			}
		}
	}()

	err := configWatcher.watcher.Add(configWatcher.filePath)
	if err != nil {
		log.Errorln(err)
	}
}

func NewConfigWatcher(filePath string, onConfigChange OnChange) (configWatcher *ConfigWatcher, err error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	configWatcher = &ConfigWatcher{watcher, onConfigChange, filePath, make(chan struct {})}
	configWatcher.Watch()
	return
}