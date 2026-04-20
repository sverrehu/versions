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
	Cache       *lrumap.LRUMap
	CommitDates *lrumap.LRUMap
	changed     bool
	mutex       sync.Mutex
}

var state *State
var stateFilename string

func init() {
	gob.Register(time.Time{})
}

func InitState(stateFile string, cacheMinutes, cacheSize int) {
	state = &State{
		Cache:       lrumap.New(cacheSize, time.Duration(cacheMinutes)*time.Minute),
		CommitDates: lrumap.New(10000, 60*24*time.Hour),
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

func PutCommitTimestamp(datasource, commitId string, timestamp time.Time) {
	state.mutex.Lock()
	defer state.mutex.Unlock()
	state.CommitDates.Put(toCommitTimestampKey(datasource, commitId), timestamp)
	state.changed = true
}

func GetCommitTimestamp(datasource, commitId string) *time.Time {
	state.mutex.Lock()
	defer state.mutex.Unlock()
	value := state.CommitDates.Get(toCommitTimestampKey(datasource, commitId))
	if value == nil {
		return nil
	}
	return new(value.(time.Time))
}

func toCommitTimestampKey(datasource, commitId string) string {
	return datasource + "/" + commitId
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
