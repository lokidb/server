package destributedstate

type state struct {
	Current map[string][]byte
	aol     appendOnlyLog
}

func New() State {
	s := new(state)
	s.aol = newAol()
	s.Current = make(map[string][]byte, 100)

	return s
}

func (s *state) applay(c ichange) {
	key := c.getKey()

	switch c.getOpration() {
	case "del":
		delete(s.Current, key)
	case "set":
		s.Current[key] = c.getValue()
	}

	s.aol.append(c)
}

func (s *state) getAol() appendOnlyLog {
	return s.aol
}

func (s *state) Merge(aol appendOnlyLog) {
	newState := new(state)
	newState.aol = newAol()
	newState.Current = make(map[string][]byte, len(s.Current))

	changes := aol.getAll()
	currentChanges := s.aol.getAll()

	index1 := 0
	index2 := 0
	for index1 < len(currentChanges) || index2 < len(changes) {
		if index1 < len(currentChanges) && index2 < len(changes) {
			c1 := currentChanges[index1]
			c2 := changes[index2]

			if c1.createdAt().UnixNano() < c2.createdAt().UnixNano() {
				newState.applay(c1)
				index1 += 1
			} else {
				newState.applay(c2)
				index2 += 1
			}
		} else if index1 < len(currentChanges) {
			c1 := currentChanges[index1]
			newState.applay(c1)
			index1 += 1
		} else {
			c2 := changes[index2]
			newState.applay(c2)
			index2 += 1
		}
	}

	s.aol = newState.aol
	s.Current = newState.Current
}

func (s *state) Del(key string) {
	s.applay(newChange("del", key, nil))
}

func (s *state) Set(key string, value []byte) {
	s.applay(newChange("set", key, value))
}

func (s *state) Get(key string) []byte {
	return s.Current[key]
}
