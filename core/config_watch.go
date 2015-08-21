package core

import (
	"log"
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
					log.Println("Modified config file: ", event.Name)
					configWatcher.onConfigChange()
				}
			case err := <-configWatcher.watcher.Errors:
				log.Println("Watch config error:", err)
			}
		}
	}()

	err := configWatcher.watcher.Add(configWatcher.filePath)
	if err != nil {
		log.Fatal(err)
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