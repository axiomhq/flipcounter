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
	entries := uint64(100)
	maxHits := uint64(10000000)
	hits := 0

	for i := uint64(0); i < entries; i++ {
		id := []byte(fmt.Sprintf("flow-%05d", i))
		expected[string(id)] = (rand.Uint64() % maxHits) + 1
		fmt.Println(string(id), expected[string(id)], hits)
		for j := uint64(0); j < expected[string(id)]; j++ {
			sk.Increment(id)
			hits++
		}
		count := sk.Query(id)
		ratio := 100*float64(count)/float64(expected[string(id)]) - 100
		if math.Abs(ratio) > 5 {
			t.Errorf("%s expected (%d != %d) ratio <= 5%%, got %2f%%", id, expected[string(id)], count, ratio)
		}
		if hits > 1000000000 {
			break
		}
	}

}
