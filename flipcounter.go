package flipcounter

import (
	"encoding/binary"
	"math"
	"math/rand"

	metro "github.com/dgryski/go-metro"
)

const (
	exp            = 1.00002
	guaranteeLimit = 8388607
	countbits      = 23
	expBase        = 255999
)

func value(c uint32) float64 {
	if c <= 1 {
		return float64(c)
	}
	return (1 - math.Pow(exp, float64(c))) / (1 - exp)
}

type key [5]byte   // guarantees us 1,099,511,627,776 unique keys
type count [3]byte // 24 bits, where 23 are counter and one is estimation flag

func hashToBytes(hash uint64) key {
	bytes := make([]byte, 10)
	_ = binary.PutUvarint(bytes, hash)
	return key{bytes[0], bytes[1], bytes[2], bytes[3]}
}

func makeKey(value []byte) key {
	return hashToBytes(metro.Hash64(value, 1337))
}

func getCount(value count) (bool, uint32) {
	p := value[2]>>7 == 1
	value[2] = value[2] << 1 >> 1
	bytes := []byte{value[0], value[1], value[2], 0}
	val := binary.LittleEndian.Uint32(bytes)
	return p, val
}

func makeCount(p bool, c uint32) count {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, c)
	if p {
		bytes[2] = bytes[2] | (1 << 7)
	}
	return count{bytes[0], bytes[1], bytes[2]}
}

// Counter consits of a 5 bytes key and a 3 bytes estimating counter with a 100% accurary up to 8388607 hits, then an estimation of 1% error
type Counter struct {
	dict map[key]count
}

// New return a HashLog sketch with 5 bytes keys and 3 bytes counters
func New() *Counter {
	return &Counter{
		dict: make(map[key]count),
	}
}

// Increment increments the counter of val []byte by +1
func (fc *Counter) Increment(val []byte) {
	k := makeKey(val)
	v := fc.dict[k]
	p, c := getCount(v)

	switch {
	case !p && c < guaranteeLimit: // good old +1 counting
		fc.dict[k] = makeCount(false, c+1)
	case !p && c == guaranteeLimit: // flip the bit, its estimation time
		fc.dict[k] = makeCount(true, expBase)
	case true && rand.Float64() < 1/math.Pow(exp, float64(c)): // roll the dice on incrementing
		fc.dict[k] = makeCount(true, c+1)
	}
}

// Get returns the number of hits for val []byte
func (fc *Counter) Get(val []byte) uint64 {
	k := makeKey(val)
	if v, ok := fc.dict[k]; ok {
		p, c := getCount(v)
		if !p {
			return uint64(c)
		}
		// if we are estimating the calculate the estimation
		return uint64(value(c))
	}
	return 0
}
