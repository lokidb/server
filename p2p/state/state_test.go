package state

import (
	"testing"
	"time"
)

func TestMergeSameID(t *testing.T) {
	s1 := New()
	s2 := New()

	s1.Update("name", "first", time.Second)
	s2.Update("name", "second", time.Second)

	s1.Merge(s2)

	items := s1.Items()

	if items[0].Value != "second" {
		t.Error("expecting item payload to be 'second'")
	}
}

func TestMergeDiffID(t *testing.T) {
	s1 := New()
	s2 := New()

	s1.Update("name1", "first", time.Second)
	s2.Update("name2", "second", time.Second)

	s1.Merge(s2)

	if len(s1.items) != 2 {
		t.Errorf("expecting 2 items in state not %d", len(s1.items))
	}
}

func TestMergeTrippleID(t *testing.T) {
	s1 := New()
	s2 := New()

	s1.Update("name", "first", time.Second)
	s2.Update("name", "second", time.Second)
	s1.Update("name", "third", time.Second)

	s1.Merge(s2)

	items := s1.Items()

	if items[0].Value != "third" {
		t.Error("expecting item payload to be 'third'")
	}
}

func TestMergeSameAndDiff(t *testing.T) {
	s1 := New()
	s2 := New()

	s1.Update("name", "first", time.Second)
	s2.Update("name", "second", time.Second)

	s1.Update("name1", "bla", time.Second)
	s2.Update("name2", "alb", time.Second)

	s1.Merge(s2)

	if len(s1.items) != 3 {
		t.Errorf("expecting 3 items after merge not %d", len(s1.items))
	}
}
