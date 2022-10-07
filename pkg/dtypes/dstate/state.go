package destributedstate

import "time"

type value struct {
	lastUpdated time.Time
	data        []byte
}

type state struct {
	Items map[string]value
}

func New() State {
	s := new(state)
	s.Items = make(map[string]value, 100)

	return s
}

func (s *state) Merge(o State) State {
	newState := new(state)
	newState.Items = make(map[string]value, len(s.Items))

	for key, value := range s.Items {
		newState.Items[key] = value
	}

	for key, value := range o.GetItems() {
		currValue, ok := newState.Items[key]
		if ok {
			if value.lastUpdated.UnixNano() > currValue.lastUpdated.UnixNano() {
				newState.Items[key] = value
			}
		} else {
			newState.Items[key] = currValue
		}
	}

	return newState
}

func (s *state) Del(key string) {
	s.Items[key] = value{lastUpdated: time.Now(), data: nil}
}

func (s *state) Set(key string, data []byte) {
	s.Items[key] = value{lastUpdated: time.Now(), data: data}
}

func (s *state) Get(key string) []byte {
	value, ok := s.Items[key]
	if !ok {
		return nil
	}

	return value.data
}

func (s *state) GetItems() map[string]value {
	return s.Items
}
