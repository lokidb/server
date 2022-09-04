package p2p

import (
	"fmt"
	"time"
)

type Desentralised interface {
	merge(other Desentralised) Desentralised
}

type message struct {
	id      string
	name    string
	payload string
	created time.Time
	maxLife time.Duration
}

func NewMessage(id string, name string, payload string, maxLife time.Duration) message {
	return message{
		id:      id,
		name:    name,
		payload: payload,
		maxLife: maxLife,
		created: time.Now(),
	}
}

func (m *message) merge(other message) (message, error) {
	if m.id != other.id {
		return message{}, fmt.Errorf("cant merge messages without the same id")
	}

	var maxLife time.Duration
	var created time.Time

	if m.created.Unix() > other.created.Unix() {
		created = other.created
		maxLife = other.maxLife
	} else {
		created = m.created
		maxLife = m.maxLife
	}

	return message{
		id:      m.id,
		payload: m.payload,
		maxLife: maxLife,
		created: created,
	}, nil
}

func (m *message) isActive() bool {
	return float64(m.created.Unix())+m.maxLife.Seconds() > float64(time.Now().Unix())
}

type State struct {
	messages map[string]message
}

func newState() State {
	return State{}
}

func (s *State) merge(other State) State {
	messages := make(map[string]message, len(s.messages)+len(other.messages))

	for id, msg := range s.messages {
		messages[id] = msg
	}

	for id, msg := range other.messages {
		messages[id] = msg
	}

	return State{messages: messages}
}

func (s *State) AddMessage(msg message) {
	s.messages[msg.id] = msg
}
