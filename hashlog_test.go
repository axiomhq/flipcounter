package hashlog

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"
)

func TestSketch(t *testing.T) {
	sk := New()
	expected := map[string]uint64{}
	entries := uint64(10)
	maxHits := uint64(100000000)

	hits := 0
	for i := uint64(0); i < entries; i++ {
		id := []byte(fmt.Sprintf("flow-%05d", i))
		fmt.Println(string(id))
		expected[string(id)] = (rand.Uint64() % maxHits) + 1
		for j := uint64(0); j < expected[string(id)]; j++ {
			sk.Increment(id)
			hits++
		}
		count := sk.Query(id)
		ratio := 100*float64(count)/float64(expected[string(id)]) - 100
		if math.Abs(ratio) > 3 {
			t.Errorf("%s expected (%d != %d) ratio <= 5%%, got %2f%% (total %d hits)", id, expected[string(id)], count, ratio, hits)
		}
	}

}

func TestSketchOverflow(t *testing.T) {
	sk := New()
	for i := 0; i < guaranteeLimit; i++ {
		sk.Increment([]byte("seif"))
	}
	count := sk.Query([]byte("seif"))
	if count != guaranteeLimit {
		t.Errorf("expected %d, got %d", guaranteeLimit, count)
	}

	sk.Increment([]byte("seif"))
	count = sk.Query([]byte("seif"))
	ratio := 100*float64(count)/float64(guaranteeLimit+1) - 100
	if math.Abs(ratio) > 1 {
		t.Errorf("%s expected (%d != %d) ratio <= 5%%, got %2f%%", string([]byte("seif")), guaranteeLimit+1, count, ratio)
	}
}

func TestValSplitJoin(t *testing.T) {
	bits := "00000110100000000000000000000101"
	parsed, _ := strconv.ParseUint(bits, 2, 32)
	val := uint32(parsed)
	k, p, c := splitVal(val)
	joined := joinVal(k, p, c)
	if val != joined {
		t.Errorf("expected %d, got %d", val, joined)
	}
}
