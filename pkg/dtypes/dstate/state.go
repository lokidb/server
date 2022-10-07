package dstate

import "time"

type Value struct {
	LastUpdated time.Time
	Data        []byte
}

type state struct {
	Items map[string]Value
}

func New() State {
	s := new(state)
	s.Items = make(map[string]Value, 100)

	return s
}

func (s *state) Merge(o State) State {
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
	s.Items[key] = Value{LastUpdated: time.Now(), Data: nil}
}

func (s *state) Set(key string, data []byte) {
	s.Items[key] = Value{LastUpdated: time.Now(), Data: data}
}

func (s *state) Get(key string) []byte {
	value, ok := s.Items[key]
	if !ok {
		return nil
	}

	return value.Data
}

func (s *state) GetItems() *map[string]Value {
	return &s.Items
}
