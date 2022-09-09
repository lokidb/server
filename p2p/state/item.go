package state

import (
	"time"
)

type item struct {
	Name             string
	Key              string
	Value            any
	LastUpdate       time.Time
	InactiveDuration time.Duration
}

func newItem(name string, key string, value any, inactiveDuration time.Duration) item {
	return item{Name: name, Key: key, Value: value, LastUpdate: time.Now().UTC()}
}

func (i *item) fullId() string {
	return i.Key + " " + i.Name
}
