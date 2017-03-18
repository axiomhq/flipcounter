package flipcounter

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestCounter(t *testing.T) {
	sk := New()
	expected := map[string]uint64{}
	entries := uint64(10)
	maxHits := uint64(100000000)

	hits := 0
	for i := uint64(0); i < entries; i++ {
		id := []byte(fmt.Sprintf("flow-%05d", i))
		expected[string(id)] = (rand.Uint64() % maxHits) + 1
		for j := uint64(0); j < expected[string(id)]; j++ {
			sk.Increment(id)
			hits++
		}
		count := sk.Get(id)
		ratio := 100*float64(count)/float64(expected[string(id)]) - 100
		if math.Abs(ratio) > 1 {
			t.Errorf("%s expected (%d != %d) ratio <= 1%%, got %2f%% (total %d hits)", id, expected[string(id)], count, ratio, hits)
			return
		}
	}

}

func TestCounterBillion(t *testing.T) {
	sk := New()
	expected := 1000000000
	for i := 0; i < expected; i++ {
		sk.Increment([]byte("seif"))
	}
	count := sk.Get([]byte("seif"))
	ratio := 100*float64(count)/float64(expected) - 100
	if math.Abs(ratio) > 1 {
		t.Errorf("expected (%d != %d) ratio <= 1%%, got %2f%%", expected, count, ratio)
		return
	}
}

func TestCounterOverflow(t *testing.T) {
	sk := New()
	for i := 0; i < guaranteeLimit; i++ {
		sk.Increment([]byte("seif"))
	}
	count := sk.Get([]byte("seif"))
	if count != guaranteeLimit {
		t.Errorf("expected %d, got %d", guaranteeLimit, count)
	}

	sk.Increment([]byte("seif"))
	count = sk.Get([]byte("seif"))
	ratio := 100*float64(count)/float64(guaranteeLimit+1) - 100
	if math.Abs(ratio) > 1 {
		t.Errorf("%s expected (%d != %d) ratio <= 1%%, got %2f%%", string([]byte("seif")), guaranteeLimit+1, count, ratio)
	}
}

func TestCounterTwice(t *testing.T) {
	sk := New()
	sk.Increment([]byte("seif"))
	count := sk.Get([]byte("seif"))
	if count != 1 {
		t.Errorf("expected %d, got %d", 1, count)
	}
	sk.Increment([]byte("seif"))
	count = sk.Get([]byte("seif"))
	if count != 2 {
		t.Errorf("expected %d, got %d", 2, count)
	}
}

func TestGetCount(t *testing.T) {
	expected := uint32(2097152)
	count := makeCount(false, expected)
	p, got := getCount(count)
	if p != false {
		t.Errorf("expected 0 == %v, got %v", false, p)
	}
	if got != expected {
		t.Errorf("expected count == %d, got %d", expected, got)
	}
}
