package p2p

import (
	"testing"
	"time"
)

func TestMessageCreation(t *testing.T) {
	msg := NewMessage("new-node 1.5.9.8", "new-node", "1.5.9.8", time.Minute*1)

	// the message creation date is less then 2 seconds ago
	if time.Now().Unix()-msg.created.Unix() > 2 {
		t.Error("expecting creation date to be now")
	}

	if !msg.isActive() {
		t.Error("expecting message to be active")
	}
}

func TestMerge(t *testing.T) {
	m1 := NewMessage("test1", "test", "testest", time.Second)
	m2 := NewMessage("test1", "test", "testest", time.Second*2)

	m3, _ := m1.merge(m2)

	if m3.id != "test1" {
		t.Error("expecting id to be as m1 id")
	}

	if m3.payload != "testest" {
		t.Error("expecting payload to be as m1 payload")
	}

	if m3.maxLife != m1.maxLife || m3.created != m1.created {
		t.Error("expecting created and max life to be as m1")
	}
}

func TestIsActive(t *testing.T) {
	m := NewMessage("a", "test", "", time.Second*1)

	if !m.isActive() {
		t.Error("expecting active message")
	}

	time.Sleep(time.Millisecond * 1020)

	if m.isActive() {
		t.Error("expecting inactive message")
	}
}

func TestStateMerge(t *testing.T) {
	s1 := newState()
	s1.AddMessage(NewMessage("a", "test", "", time.Second))

	s2 := newState()
	s2.AddMessage(NewMessage("a", "test", "", time.Second*5))

	s3 := s1.merge(s2)

	if len(s3.messages) != 1 {
		t.Fatal("expecting one message")
	}
}

func TestAddMessage(t *testing.T) {
	s := newState()
	s.AddMessage(NewMessage("a", "test", "", time.Second))

	if len(s.messages) != 1 {
		t.Error("expecting one message")
	}

	s.AddMessage(NewMessage("a", "test", "", time.Second))

	if len(s.messages) != 1 {
		t.Error("expecting one message")
	}

	s.AddMessage(NewMessage("b", "test", "", time.Second))

	if len(s.messages) != 2 {
		t.Error("expecting 2 messages")
	}
}

func TestRemoveInactiveMessages(t *testing.T) {
	s1 := newState()
	s2 := newState()
	s1.AddMessage(NewMessage("a", "test", "", time.Second))

	time.Sleep(time.Second)

	s3 := s1.merge(s2)

	if len(s3.messages) != 0 {
		t.Error("expecting zero messages")
	}
}
