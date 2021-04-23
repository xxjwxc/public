package kv

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/tidwall/gjson"
)

func newWatcher(path string) (*watcher, error) {
	wp, err := watch.Parse(map[string]interface{}{"type": "keyprefix", "prefix": path})
	if err != nil {
		return nil, err
	}

	return &watcher{
		Plan:       wp,
		lastValues: make(map[string][]byte),
		err:        make(chan error, 1),
	}, nil
}

type watcher struct {
	*watch.Plan
	lastValues    map[string][]byte
	hybridHandler watch.HybridHandlerFunc
	stopChan      chan interface{}
	err           chan error
	sync.RWMutex
}

func (w *watcher) getValue(path string) []byte {
	w.RLock()
	defer w.RUnlock()

	return w.lastValues[path]
}

func (w *watcher) updateValue(path string, value []byte) {
	w.Lock()
	defer w.Unlock()

	if len(value) == 0 {
		delete(w.lastValues, path)
	} else {
		w.lastValues[path] = value
	}
}

func (w *watcher) setHybridHandler(prefix string, handler func(*Result)) {
	w.hybridHandler = func(bp watch.BlockingParamVal, data interface{}) {
		kvPairs := data.(api.KVPairs)
		ret := &Result{}

		for _, k := range kvPairs {
			path := strings.TrimSuffix(strings.TrimPrefix(k.Key, prefix+"/"), "/")
			v := w.getValue(path)

			if len(k.Value) == 0 && len(v) == 0 {
				continue
			}

			if bytes.Equal(k.Value, v) {
				continue
			}

			ret.g = gjson.ParseBytes(k.Value)
			ret.k = path
			w.updateValue(path, k.Value)
			handler(ret)
		}
	}
}

func (w *watcher) run(address string, conf *api.Config) error {
	w.stopChan = make(chan interface{})
	w.Plan.HybridHandler = w.hybridHandler

	go func() {
		w.err <- w.RunWithConfig(address, conf)
	}()

	select {
	case err := <-w.err:
		return fmt.Errorf("run fail: %w", err)
	case <-w.stopChan:
		w.Stop()
		return nil
	}
}

func (w *watcher) stop() {
	close(w.stopChan)
}
