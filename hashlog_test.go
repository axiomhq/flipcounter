package hashlog

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestSketch(t *testing.T) {
	sk := New()
	expected := map[string]uint64{}
	entries := rand.Uint64() % 1e4
	maxHits := uint64(1e5)

	for i := uint64(0); i < entries; i++ {
		id := []byte(fmt.Sprintf("flow-%05d", i))
		expected[string(id)] = rand.Uint64() % maxHits
		for j := uint64(0); j < expected[string(id)]; j++ {
			sk.Increment(id)
		}
		ratio := 100*float64(sk.Query(id))/float64(expected[string(id)]) - 100
		if math.Abs(ratio) >= 5 {
			t.Errorf("expected ratio < 5%%, got %2f%%", ratio)
		}
	}

}
