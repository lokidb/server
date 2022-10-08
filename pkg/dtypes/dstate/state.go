package dstate

import (
	"sync"
	"time"
)

type Value struct {
	LastUpdated time.Time
	Data        []byte
}

type state struct {
	Items map[string]Value
	lock  sync.RWMutex
}

func New() State {
	s := new(state)
	s.Items = make(map[string]Value, 100)

	return s
}

func (s *state) Merge(o State) State {
	s.lock.RLock()
	defer s.lock.RUnlock()

	newState := new(state)
	newState.Items = make(map[string]Value, len(s.Items))

	for key, value := range s.Items {
		newState.Items[key] = value
	}

	for key, value := range *o.GetItems() {
		currValue, ok := newState.Items[key]
		if ok {
			if value.LastUpdated.UnixNano() > currValue.LastUpdated.UnixNano() {
				newState.Items[key] = value
			}
		} else {
			newState.Items[key] = currValue
		}
	}

	return newState
}

func (s *state) Del(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.Items[key] = Value{LastUpdated: time.Now(), Data: nil}
}

func (s *state) Set(key string, data []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Items[key] = Value{LastUpdated: time.Now(), Data: data}
}

func (s *state) SetDefault(key string, data []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Items[key] = Value{LastUpdated: time.Unix(0, 0), Data: data}
}

func (s *state) Get(key string) []byte {
	s.lock.RLock()
	defer s.lock.RUnlock()

	value, ok := s.Items[key]
	if !ok {
		return nil
	}

	return value.Data
}

func (s *state) GetItems() *map[string]Value {
	return &s.Items
}
