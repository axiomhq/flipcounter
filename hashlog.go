package hashlog

import (
	"math"
	"math/rand"

	metro "github.com/dgryski/go-metro"
)

const (
	exp             = 1.00026
	keybits         = 32
	fingerprintbits = 8
	precisionbits   = 1
	countbits       = 64 - keybits - fingerprintbits - precisionbits
	guaranteeLimit  = 8388607
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

func indicies(val []byte) (uint32, uint8) {
	hash := metro.Hash64(val, 1337)
	index := uint32(hash >> keybits)
	fingerprint := uint8(hash << keybits >> (keybits + precisionbits + countbits))
	return index, fingerprint
}

func splitVal(val uint32) (uint8, bool, uint32) {
	fingerprint := uint8(val >> (countbits + precisionbits))
	precLog := uint16(val<<fingerprintbits>>(fingerprintbits+countbits)) == 1
	count := uint32(val << (precisionbits + fingerprintbits) >> (precisionbits + fingerprintbits))
	return fingerprint, precLog, count
}

func joinVal(key uint8, precLog bool, count uint32) uint32 {
	val := uint32(key) << (countbits + precisionbits)
	if precLog {
		val |= (1 << countbits)
	}
	val |= uint32(count)
	return val
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

func (sketch *Sketch) get(h1 uint32, h2 uint8) uint64 {
	for _, val := range sketch.dict[h1] {
		if key, precLog, count := splitVal(val); key == h2 {
			if !precLog {
				return uint64(count)
			}
			return uint64(value(count))
		}
	}
	return 0
}

func (sketch *Sketch) inc(h1 uint32, h2 uint8) {
	for i, val := range sketch.dict[h1] {
		key, precLog, count := splitVal(val)
		if key == h2 {
			if !precLog {
				if count < guaranteeLimit {
					sketch.dict[h1][i] = joinVal(key, false, count+1)
				} else {
					sketch.dict[h1][i] = joinVal(key, true, 29574)
				}
			} else {
				if rand.Float64() < 1/math.Pow(exp, float64(count)) {
					sketch.dict[h1][i] = joinVal(key, true, count+1)
				}
			}
			return
		}
	}
	sketch.dict[h1] = append(sketch.dict[h1], joinVal(h2, false, 1))
}

// Increment ...
func (sketch *Sketch) Increment(val []byte) {
	h1, h2 := indicies(val)
	sketch.inc(h1, h2)
}

// Query ...
func (sketch *Sketch) Query(val []byte) uint64 {
	h1, h2 := indicies(val)
	return sketch.get(h1, h2)
}
