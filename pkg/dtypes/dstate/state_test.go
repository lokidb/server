package destributedstate

import (
	"testing"
)

func equal(a []byte, b []byte) bool {
	if a == nil && b == nil {
		return true
	}

	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestState(t *testing.T) {
	s := New()

	s.Set("a", []byte("abc"))
	s.Set("b", []byte("yty"))
	s.Del("a")

	if s.Get("a") != nil {
		t.Errorf("expecting nil on deleted key")
	}

	if !equal(s.Get("b"), []byte("yty")) {
		t.Errorf("expecting value 'yty' for key 'b'")
	}
}

func TestMerge(t *testing.T) {
	s1 := New()
	s2 := New()

	s1.Set("a", []byte("A"))
	s1.Set("c", []byte("C"))
	s1.Set("a", []byte("B"))
	s2.Set("a", []byte("A"))
	s1.Set("b", []byte("B"))
	s1.Del("a")
	s2.Set("a", []byte("A"))
	s2.Del("c")

	s3 := s1.Merge(s2)

	if !equal(s3.Get("a"), []byte("A")) {
		t.Errorf("expecting 'a' value to be 'A'")
	}

	if !equal(s3.Get("b"), []byte("B")) {
		t.Errorf("expecting 'b' value to be 'B'")
	}

	if !equal(s3.Get("c"), nil) {
		t.Errorf("expecting 'c' value to be nil")
	}
}
