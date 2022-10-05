package destributedstate

import "time"

type appendOnlyLog interface {
	append(ichange)
	getAll() []ichange
}

type ichange interface {
	getOpration() string
	getKey() string
	getValue() []byte
	createdAt() time.Time
}

type State interface {
	Merge(appendOnlyLog)
	Set(key string, value []byte)
	Get(key string) []byte
	Del(key string)
	applay(ichange)
	getAol() appendOnlyLog
}
