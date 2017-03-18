package hashlog

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
	switch c {
	case 0:
		return 0
	case 1:
		return math.Pow(exp, float64(c-1))
	default:
		return (1 - math.Pow(exp, float64(c))) / (1 - exp)
	}
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

// Sketch consits of a 5 bytes key and a 3 bytes estimating counter with a 100% accurary up to 8388607 hits, then an estimation of 1% error
type Sketch struct {
	dict map[key]count
}

// New return a HashLog sketch with 5 bytes keys and 3 bytes counters
func New() *Sketch {
	return &Sketch{
		dict: make(map[key]count),
	}
}

// Increment increments the counter of val []byte by +1
func (sketch *Sketch) Increment(val []byte) {
	k := makeKey(val)
	v := sketch.dict[k]
	p, c := getCount(v)

	switch {
	case !p && c < guaranteeLimit:
		sketch.dict[k] = makeCount(false, c+1)
	case !p && c == guaranteeLimit:
		sketch.dict[k] = makeCount(true, expBase)
	case true && rand.Float64() < 1/math.Pow(exp, float64(c)):
		sketch.dict[k] = makeCount(true, c+1)
	}
}

// Count returns the number of hits for val []byte
func (sketch *Sketch) Count(val []byte) uint64 {
	k := makeKey(val)
	v, ok := sketch.dict[k]
	if !ok {
		return 0
	}
	p, c := getCount(v)
	if !p {
		return uint64(c)
	}
	// if we are estimating the calculate the estimation
	return uint64(value(c))
}
