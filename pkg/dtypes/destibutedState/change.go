package destributedstate

import "time"

type change struct {
	operation string
	key       string
	value     []byte
	created   time.Time
}

func newChange(operation string, key string, value []byte) ichange {
	return &change{operation: operation, key: key, value: value, created: time.Now()}
}

func (c *change) createdAt() time.Time {
	return c.created
}

func (c *change) getKey() string {
	return c.key
}

func (c *change) getValue() []byte {
	return c.value
}

func (c *change) getOpration() string {
	return c.operation
}
