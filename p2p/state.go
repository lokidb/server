package p2p

import (
	"fmt"
	"time"
)

type Desentralised interface {
	merge(other Desentralised) Desentralised
	sign(signature string)
}

type message struct {
	id        string
	name      string
	payload   string
	created   time.Time
	maxLife   time.Duration
	signaturs []string
	handled   bool
}

func NewMessage(id string, name string, payload string, maxLife time.Duration) message {
	return message{
		id:        id,
		name:      name,
		payload:   payload,
		maxLife:   maxLife,
		created:   time.Now(),
		signaturs: make([]string, 0, 20),
		handled:   false,
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

	// Merge the signatures lists without duplicates
	signaturs := make([]string, 0, len(m.signaturs)+len(other.signaturs))
	signaturs = append(signaturs, m.signaturs...)

	for _, s := range other.signaturs {
		exists := false
		for i := range signaturs {
			if signaturs[i] == s {
				exists = true
			}
		}

		if !exists {
			signaturs = append(signaturs, s)
		}
	}

	return message{
		id:        m.id,
		payload:   m.payload,
		maxLife:   maxLife,
		created:   created,
		signaturs: signaturs,
	}, nil
}

func (m *message) isActive() bool {
	return float64(m.created.Unix())+m.maxLife.Seconds() > float64(time.Now().Unix())
}

func (m *message) sign(signature string) {
	for _, sig := range m.signaturs {
		if sig == signature {
			return
		}
	}

	m.signaturs = append(m.signaturs, signature)
}

type State struct {
	messages []message
}

func newState() State {
	return State{}
}

func (s *State) merge(other State) State {
	messages := make([]message, 0, len(s.messages)+len(other.messages))

	for _, msg := range s.messages {
		if msg.isActive() {
			messages = append(messages, msg)
		}
	}

	for _, msg := range other.messages {
		if !msg.isActive() {
			continue
		}

		replaced := false
		for i, omsg := range messages {
			if msg.id == omsg.id {
				mmsg, _ := msg.merge(omsg)
				replaced = true
				messages[i] = mmsg
			}
		}

		if !replaced {
			messages = append(messages, msg)
		}
	}

	return State{messages: messages}
}

func (s *State) AddMessage(msg message) {
	for i, omsg := range s.messages {
		if msg.id == omsg.id {
			mmsg, _ := msg.merge(omsg)
			s.messages[i] = mmsg
			return
		}
	}

	s.messages = append(s.messages, msg)
}

func (s *State) sign(signature string) {
	for _, msg := range s.messages {
		msg.sign(signature)
	}
}
