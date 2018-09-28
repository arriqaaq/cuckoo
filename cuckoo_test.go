package cuckoo

import (
	"encoding/binary"
	"fmt"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func timeGC() time.Duration {
	start := time.Now()
	runtime.GC()
	return time.Since(start)
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
	fmt.Println()
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func TestBasic(t *testing.T) {
	f := NewCuckooFilter(20, 0.001)
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
	f := NewCuckooFilter(20, 0.001)
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

func TestCuckooInsertTime(t *testing.T) {
	f := NewCuckooFilter(uint(5000000), 0.001)

	data := make([][]byte, 5000000)
	for i := 0; i < 5000000; i++ {
		data[i] = []byte(strconv.Itoa(i))
	}

	now := time.Now()

	for i := 0; i < 5000000; i++ {
		f.Insert(data[i])
	}
	PrintMemUsage()
	fmt.Println("GC took: ", 5000000, time.Since(now))
}

func TestCuckooDeleteTime(t *testing.T) {
	f := NewCuckooFilter(uint(5000000), 0.001)

	data := make([][]byte, 5000000)
	for i := 0; i < 5000000; i++ {
		data[i] = []byte(strconv.Itoa(i))
	}

	for i := 0; i < 5000000; i++ {
		f.Insert(data[i])
	}

	now := time.Now()

	for i := 0; i < 5000000; i++ {
		f.Delete(data[i])
	}
	PrintMemUsage()
	fmt.Println("GC took: ", 5000000, time.Since(now))
}

func TestCuckooLookupTime(t *testing.T) {
	f := NewCuckooFilter(uint(5000000), 0.001)

	data := make([][]byte, 5000000)
	for i := 0; i < 5000000; i++ {
		data[i] = []byte(strconv.Itoa(i))
	}

	for i := 0; i < 5000000; i++ {
		f.Insert(data[i])
	}

	now := time.Now()

	for i := 0; i < 5000000; i++ {
		f.Lookup(data[i])
	}
	PrintMemUsage()
	fmt.Println("GC took: ", 5000000, time.Since(now))
}

func BenchmarkCuckooInsert(b *testing.B) {
	b.StopTimer()
	f := NewCuckooFilter(uint(b.N), 0.001)

	data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = []byte(strconv.Itoa(i))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		f.Insert(data[i])
	}
	PrintMemUsage()
	runtime.GC()
	fmt.Println("GC took: ", b.N, timeGC())
}

func BenchmarkCuckooLookup(b *testing.B) {
	b.StopTimer()
	f := NewCuckooFilter(uint(b.N), 0.001)
	data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = []byte(strconv.Itoa(i))
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		f.Lookup(data[n])
	}
	PrintMemUsage()
	runtime.GC()
	fmt.Println("GC took: ", b.N, timeGC())
}

func BenchmarkCuckooLookupAndInsert(b *testing.B) {
	f := NewCuckooFilter(uint(b.N), 0.001)
	key := make([]byte, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		binary.BigEndian.PutUint32(key, uint32(i))
		f.Lookup(key)
		f.Insert(key)
	}
	PrintMemUsage()
	runtime.GC()
	fmt.Println("GC took: ", b.N, timeGC())
}

func BenchmarkCuckooLookupAndDelete(b *testing.B) {
	b.StopTimer()
	f := NewCuckooFilter(uint(b.N), 0.001)
	data := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		data[i] = []byte(strconv.Itoa(i))
	}
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		f.Insert(data[n])
		f.Delete(data[n])
	}
	PrintMemUsage()
	runtime.GC()
	fmt.Println("GC took: ", b.N, timeGC())
}
