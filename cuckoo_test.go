package cuckoo

import (
	"encoding/binary"
	"testing"
)

func TestBasic(t *testing.T) {
	f := NewCuckooFilter(20)
	n1 := []byte("Bess")
	n2 := []byte("Jane")
	f.Insert(n1)
	n1b := f.Lookup(n1)
	n2b := f.Lookup(n2)
	if !n1b {
		t.Errorf("%v should be in.", n1)
	}
	if n2b {
		t.Errorf("%v should not be in.", n2)
	}
}

func TestBasicUint32(t *testing.T) {
	f := NewCuckooFilter(20)
	n1 := make([]byte, 4)
	n2 := make([]byte, 4)
	n3 := make([]byte, 4)
	n4 := make([]byte, 4)
	binary.BigEndian.PutUint32(n1, 100)
	binary.BigEndian.PutUint32(n2, 101)
	binary.BigEndian.PutUint32(n3, 102)
	binary.BigEndian.PutUint32(n4, 103)
	f.Insert(n1)
	n1b := f.Lookup(n1)
	n2b := f.Lookup(n2)
	n3b := f.Lookup(n3)
	f.Lookup(n4)
	if !n1b {
		t.Errorf("%v should be in.", n1)
	}
	if n2b {
		t.Errorf("%v should not be in.", n2)
	}
	if n3b {
		t.Errorf("%v should not be in.", n3)
	}
}

func BenchmarkSeparateLookup(b *testing.B) {
	f := NewCuckooFilter(uint(b.N))
	b.ResetTimer()
	n1 := []byte("Bess")
	f.Insert(n1)

	for i := 0; i < b.N; i++ {
		f.Lookup(n1)
	}
}

func BenchmarkSeparateInsert(b *testing.B) {
	f := NewCuckooFilter(uint(b.N))
	b.ResetTimer()

	n1 := []byte("Bess")

	for i := 0; i < b.N; i++ {
		f.Insert(n1)
	}
}

func BenchmarkSeparateLookupAndInsert(b *testing.B) {
	f := NewCuckooFilter(uint(b.N))
	key := make([]byte, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		binary.BigEndian.PutUint32(key, uint32(i))
		f.Lookup(key)
		f.Insert(key)
	}
}
