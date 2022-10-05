package state

import (
	"time"
)

type State struct {
	items map[string]Item
}

func New() State {
	s := new(State)
	s.items = make(map[string]Item, 1000)

	return *s
}

func (s *State) Merge(other State) {
	for _, oitem := range other.items {
		item, ok := s.items[oitem.UID()]

		if ok && item.LastUpdate.After(oitem.LastUpdate) {
			continue
		}

		s.items[oitem.UID()] = oitem
	}
}

func (s *State) Update(key string, value string, inactiveDuration time.Duration) {
	i := newItem(key, value, inactiveDuration)
	s.items[i.UID()] = i
}

func (s *State) Get(key string) string {
	val, ok := s.items[key]
	if !ok {
		return ""
	}

	return val.Value
}

func (s *State) Items() []Item {
	items := make([]Item, 0, len(s.items))

	for _, item := range s.items {
		items = append(items, item)
	}

	return items
}

func (s *State) ActiveItems() []Item {
	items := make([]Item, 0, len(s.items))
	now := time.Now().UTC()

	for _, item := range s.items {
		if now.UnixMilli()-item.LastUpdate.UnixMilli() > item.InactiveDuration.Milliseconds() {
			items = append(items, item)
		}
	}

	return items
}
