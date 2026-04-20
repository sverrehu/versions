package state

import (
	"sync"
	"time"

	"github.com/sverrehu/goutils/lrumap"
)

type State struct {
	Cache   *lrumap.LRUMap
	changed bool
	mutex   sync.Mutex
}

var state = State{}

func InitState(stateFile string, cacheMinutes, cacheSize int) {
	state.Cache = lrumap.New(cacheSize, time.Duration(cacheMinutes)*time.Minute)
}

func PutCachedResponse(path string, data []byte) {
	state.Cache.Put(path, data)
	state.changed = true
}

func GetCachedResponse(path string) []byte {
	value := state.Cache.Get(path)
	if value == nil {
		return nil
	}
	return value.([]byte)
}
