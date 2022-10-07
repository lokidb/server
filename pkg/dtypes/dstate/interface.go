package destributedstate

type State interface {
	Merge(State) State
	GetItems() map[string]value
	Set(key string, value []byte)
	Get(key string) []byte
	Del(key string)
}
