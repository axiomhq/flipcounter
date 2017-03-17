package hashlog

import (
	"math"
	"math/rand"

	metro "github.com/dgryski/go-metro"
)

const exp = 1.00026

func value(c uint16) float64 {
	switch c {
	case 0:
		return 0
	case 1:
		return math.Pow(exp, float64(c-1))
	default:
		return (1 - math.Pow(exp, float64(c))) / (1 - exp)
	}
}

func indicies(val []byte) (uint32, uint16) {
	hash := metro.Hash64(val, 1337)
	return uint32(hash >> 32), uint16(hash << 48 >> 48)
}

func splitVal(val uint32) (uint16, uint16) {
	return uint16(val >> 16), uint16(val << 16 >> 16)
}

func joinVal(key, count uint16) uint32 {
	return uint32(key)<<16 | uint32(count)
}

// Sketch ...
type Sketch struct {
	dict map[uint32][]uint32
}

// New ...
func New() *Sketch {
	return &Sketch{
		dict: make(map[uint32][]uint32),
	}
}

func (sketch *Sketch) get(h1 uint32, h2 uint16) uint16 {
	for _, val := range sketch.dict[h1] {
		if key, count := splitVal(val); key == h2 {
			return count
		}
	}
	return 0
}

func (sketch *Sketch) inc(h1 uint32, h2 uint16) {
	for i, val := range sketch.dict[h1] {
		key, count := splitVal(val)
		if key == h2 {
			if rand.Float64() < 1/math.Pow(exp, float64(count)) {
				sketch.dict[h1][i] = joinVal(key, count+1)
			}
			return
		}
	}
	sketch.dict[h1] = append(sketch.dict[h1], joinVal(h2, 1))
}

// Increment ...
func (sketch *Sketch) Increment(val []byte) {
	h1, h2 := indicies(val)
	sketch.inc(h1, h2)
}

// Query ...
func (sketch *Sketch) Query(val []byte) uint64 {
	h1, h2 := indicies(val)
	return uint64(value(sketch.get(h1, h2)))
}
