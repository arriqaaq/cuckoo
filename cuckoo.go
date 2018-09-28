package cuckoo

import (
	"bytes"
	"encoding/binary"
	"errors"
	// "fmt"
	"github.com/cespare/xxhash"
	"hash/fnv"
	"math"
	"math/rand"
)

var (
	ERR_FILTER_FULL = errors.New("full")
)

func calculateFingerprintSize(b uint, epsilon float64) uint {
	f := uint(math.Ceil(math.Log(2 * float64(b) / epsilon)))
	f = f / 8
	if f <= 0 {
		f = 1
	}
	return f
}

func max(x, y uint) uint {
	if x > y {
		return x
	}
	return y
}

func hash1(data []byte) []byte {
	h := fnv.New32()
	h.Write(data)
	return h.Sum(nil)
}

func hash2(data []byte) []byte {
	h := xxhash.New()
	h.Write(data)
	return h.Sum(nil)
}

/*
	This procedure may repeat until a vacant bucket
	is found as illustrated in Figure 1(b), or until
	a maximum number of displacements is reached
	(e.g., 500 times in our implementation).
*/

const maxCuckooKicks = 500
const bucketSize = uint(4)

/*
	Fig 1(c)
	A cuckoo filter, two hash per item and
	functions and four entries per bucket
*/

type fingerprint []byte

func newBucket(size uint) bucket {
	return make(bucket, size)
}

type bucket []fingerprint

func (b bucket) contains(f fingerprint) bool {
	return b.indexOf(f) != -1
}

func (b bucket) indexOf(f fingerprint) int {
	for i, item := range b {
		if bytes.Equal(f, item) {
			return i
		}
	}
	return -1
}

func (b bucket) delete(f fingerprint) bool {
	for i, item := range b {
		if bytes.Equal(f, item) {
			b[i] = nil
			return true
		}
	}
	return false
}

func (b bucket) getEmptyEntry() (int, error) {
	for i, item := range b {
		if item == nil {
			return i, nil
		}
	}
	return -1, ERR_FILTER_FULL
}

/*
	Fig 1(c)
	A cuckoo filter, two hash per item and
	functions and four entries per bucket
*/

func NewCuckooFilter(capacity uint, fpRate float64) *CuckooFilter {
	capacity = power2(capacity) / bucketSize
	if capacity == 0 {
		capacity = 1
	}
	buckets := make([]bucket, capacity)
	for i := range buckets {
		buckets[i] = newBucket(bucketSize)
	}

	return &CuckooFilter{
		fpSize:         calculateFingerprintSize(bucketSize, fpRate),
		numbuckets:     capacity,
		buckets:        buckets,
		entryPerBucket: uint(4),
	}
}

type CuckooFilter struct {
	numbuckets     uint
	entryPerBucket uint
	fpSize         uint
	size           uint // number of items in filter
	buckets        []bucket
}

func (c *CuckooFilter) getCuckooParams(data []byte) (uint, uint, []byte) {
	hash := c.computeHash(data)
	f := hash[0:c.fpSize]
	i1 := uint(binary.BigEndian.Uint32(hash2(hash))) % c.numbuckets
	i2 := (i1 ^ uint(binary.BigEndian.Uint32(hash1(f)))) % c.numbuckets

	return i1, i2, f
}

func (c *CuckooFilter) computeHash(data []byte) []byte {
	return hash1(data)
}

func (c *CuckooFilter) insert(i1 uint, i2 uint, f []byte) error {
	// fmt.Println("vals: ", i1, i2, f, len(c.buckets))
	// insert into bucket1
	b1 := c.buckets[i1]
	if idx, err := b1.getEmptyEntry(); err == nil {
		b1[idx] = f
		return nil
	}

	// insert into bucket2
	b2 := c.buckets[i2]
	if idx, err := b2.getEmptyEntry(); err == nil {
		b2[idx] = f
		return nil
	}

	// must relocate existing items

	// i = randomly pick i1 or i2;
	// for n = 0; n < MaxNumKicks; n++ do
	// randomly select an entry e from bucket[i];
	// swap f and the fingerprint stored in entry e;
	// i = i ⊕ hash(f);
	// if bucket[i] has an empty entry then
	// add f to bucket[i]

	i := i1
	for n := 0; n < maxCuckooKicks; n++ {
		// randomly select an entry e from bucket[i];
		rIdx := rand.Intn(len(c.buckets[i]) - 1)
		f, c.buckets[i][rIdx] = c.buckets[i][rIdx], f
		i = (i ^ uint(binary.BigEndian.Uint32(hash2(f)))) % c.numbuckets
		b := c.buckets[i]
		if idx, err := b.getEmptyEntry(); err == nil {
			b[idx] = f
			return nil
		}
	}

	return ERR_FILTER_FULL

}

func (c *CuckooFilter) Insert(data []byte) error {
	i1, i2, f := c.getCuckooParams(data)
	return c.insert(i1, i2, f)
}

func (c *CuckooFilter) Lookup(data []byte) bool {
	// f = fingerprint(x);
	// i1 = hash(x);
	// i2 = i1 ⊕ hash(f);
	// if bucket[i1] or bucket[i2] has f then
	// return True;
	// return False;

	i1, i2, f := c.getCuckooParams(data)
	if c.buckets[i1].contains(f) || c.buckets[i2].contains(f) {
		return true
	}
	return false
}

func (c *CuckooFilter) Delete(data []byte) bool {
	i1, i2, f := c.getCuckooParams(data)
	if c.buckets[i1].delete(f) || c.buckets[i2].delete(f) {
		return true
	}
	return false
}

func power2(x uint) uint {
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	x |= x >> 32
	x++
	return x
}
