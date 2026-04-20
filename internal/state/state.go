package state

import (
	"encoding/gob"
	"log"
	"os"
	"sync"
	"time"

	"github.com/sverrehu/goutils/lrumap"
)

type State struct {
	Cache   *lrumap.LRUMap
	changed bool
	mutex   sync.Mutex
}

var state *State
var stateFilename string

func InitState(stateFile string, cacheMinutes, cacheSize int) {
	state = &State{
		Cache: lrumap.New(cacheSize, time.Duration(cacheMinutes)*time.Minute),
	}
	stateFilename = stateFile
	if len(stateFilename) != 0 {
		err := LoadState()
		if err != nil {
			log.Printf("error loading state, using empty state: %v", err)
		}
		go periodicStateSaveTask()
	}
}

func SaveState() error {
	if len(stateFilename) == 0 || !state.changed {
		return nil
	}
	log.Printf("saving state to %v", stateFilename)
	fh, err := os.Create(stateFilename)
	if err != nil {
		return err
	}
	defer fh.Close()
	state.mutex.Lock()
	defer state.mutex.Unlock()
	err = gob.NewEncoder(fh).Encode(&state)
	state.changed = false
	return err
}

func LoadState() error {
	if len(stateFilename) == 0 {
		return nil
	}
	log.Printf("loading state from %s", stateFilename)
	fh, err := os.Open(stateFilename)
	if err != nil {
		return err
	}
	defer fh.Close()
	newState := State{}
	err = gob.NewDecoder(fh).Decode(&newState)
	if err != nil {
		return err
	}
	state.mutex.Lock()
	defer state.mutex.Unlock()
	newState.changed = false
	state = &newState
	return nil
}

func PutCachedResponse(path string, data []byte) {
	state.mutex.Lock()
	defer state.mutex.Unlock()
	state.Cache.Put(path, data)
	state.changed = true
}

func GetCachedResponse(path string) []byte {
	state.mutex.Lock()
	defer state.mutex.Unlock()
	value := state.Cache.Get(path)
	if value == nil {
		return nil
	}
	return value.([]byte)
}

func periodicStateSaveTask() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		err := SaveState()
		if err != nil {
			log.Printf("error saving state, ignoring: %v", err)
		}
	}
}
