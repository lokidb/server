package dstate

type State interface {
	Merge(State) State
	GetItems() *map[string]Value
	Set(key string, value []byte)
	Get(key string) []byte
	Del(key string)
}
