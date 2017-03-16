package hashlog

import (
	"math"

	metro "github.com/dgryski/go-metro"
	"github.com/dgryski/go-pcgr"
)

const exp = 1.00026

var rnd = pcgr.Rand{State: 0x0ddc0ffeebadf00d, Inc: 0xcafebabe}

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

// Sketch ...
type Sketch struct {
	dict map[uint64]uint16
}

// New ...
func New() *Sketch {
	return &Sketch{
		dict: make(map[uint64]uint16),
	}
}

// Increment ...
func (sketch *Sketch) Increment(val []byte) {
	hash := metro.Hash64(val, 1337)
	c := sketch.dict[hash]
	if float64(rnd.Next()%10e5)/10e5 < 1/math.Pow(exp, float64(c)) {
		sketch.dict[hash] = sketch.dict[hash] + 1
	}
}

// Query ...
func (sketch *Sketch) Query(val []byte) uint64 {
	return uint64(value(sketch.dict[metro.Hash64(val, 1337)]))
}
