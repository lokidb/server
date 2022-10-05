package state

import (
	"time"
)

type Item struct {
	Name             string
	Key              string
	Value            string
	LastUpdate       time.Time
	InactiveDuration time.Duration
}

func newItem(key string, value string, inactiveDuration time.Duration) Item {
	return Item{Key: key, Value: value, LastUpdate: time.Now().UTC()}
}

func (i *Item) UID() string {
	return i.Key + " " + i.Name
}
