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
	// NOTE! This struct is gob persisted; do not change the exported names!
	Cache            *lrumap.LRUMap
	CommitTimestamps *lrumap.LRUMap
	changed          bool
	mutex            sync.Mutex
}

var state *State
var stateFilename string

func init() {
	gob.Register(time.Time{})
}

func InitState(stateFile string, responseCacheMinutes, responseCacheSize, commitTimestampCacheMinutes, commitTimestampCacheSize int) {
	state = &State{
		Cache:            lrumap.New(responseCacheSize, time.Duration(responseCacheMinutes)*time.Minute),
		CommitTimestamps: lrumap.New(commitTimestampCacheSize, time.Duration(commitTimestampCacheMinutes)*time.Minute),
	}
	stateFilename = stateFile
	if len(stateFilename) != 0 {
		log.Printf("state will be persisted to: %s", stateFilename)
		err := LoadState()
		if err != nil {
			log.Printf("error loading state, using empty state: %v", err)
		}
		// reset cache settings, in case current settings differ from what was persisted
		state.Cache.MaxSize = responseCacheSize
		state.Cache.TTL = time.Duration(responseCacheMinutes) * time.Minute
		state.CommitTimestamps.MaxSize = commitTimestampCacheSize
		state.CommitTimestamps.TTL = time.Duration(commitTimestampCacheMinutes) * time.Minute
		go periodicStateSaveTask()
	}
	log.Printf("state store set up with response cache minutes: %d, size: %d, and commit timestamp cache minutes: %d, size: %d",
		int(state.Cache.TTL.Minutes()), state.Cache.MaxSize, int(state.CommitTimestamps.TTL.Minutes()), state.CommitTimestamps.MaxSize)
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
	state.CommitTimestamps.Put(toCommitTimestampKey(datasource, commitId), timestamp)
	state.changed = true
}

func GetCommitTimestamp(datasource, commitId string) *time.Time {
	state.mutex.Lock()
	defer state.mutex.Unlock()
	value := state.CommitTimestamps.Get(toCommitTimestampKey(datasource, commitId))
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
