package state

import "time"

type State struct {
	items map[string]item
}

func New() State {
	s := new(State)
	s.items = make(map[string]item, 1000)

	return *s
}

func (s *State) Merge(other State) {
	for _, oitem := range other.items {
		item, ok := s.items[oitem.fullId()]

		if ok && item.LastUpdate.After(oitem.LastUpdate) {
			continue
		}

		s.items[oitem.fullId()] = oitem
	}
}

func (s *State) Update(name string, key string, value any, inactiveDuration time.Duration) {
	i := newItem(name, key, value, inactiveDuration)
	s.items[i.fullId()] = i
}

func (s *State) GetItemsByName(name string) []item {
	items := make([]item, 0, len(s.items))

	for _, item := range s.items {
		if item.Name == name {
			items = append(items, item)
		}
	}

	return items
}

func (s *State) Items() []item {
	items := make([]item, 0, len(s.items))

	for _, item := range s.items {
		items = append(items, item)
	}

	return items
}
