package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/seiflotfy/hashlog"
)

func main() {
	hl := hashlog.New()

	max := 100000
	now := time.Now()
	expected := make([]uint, max, max)
	zipf := rand.NewZipf(rand.New(rand.NewSource(now.UnixNano())), 1.1, 4.0, uint64(max)-1)

	seen := map[string]bool{}
	for k := uint64(0); len(seen) != max; k++ {
		if k%uint64(max) == 0 {
			fmt.Printf("\rCardinality %06d\t Hits: %d", len(seen), k)
		}
		i := zipf.Uint64()
		expected[i]++
		id := []byte(fmt.Sprintf("flow-%05d", i))
		seen[string(id)] = true
		hl.Increment(id)
	}

	for i := range expected {
		// some minor print for easier visuals
		if i == 100 || i-50 == (len(expected)/2)+1 || i == len(expected)-100 {
			fmt.Printf("\n---")
		}

		if (i > len(expected)-100 || i < 100) || (i < (len(expected)/2)+50 && i > (len(expected)/2)-50) {
			// id
			id := fmt.Sprintf("flow-%05d", i)
			// estimation
			est1 := float64(hl.Count([]byte(id)))
			// error ratio
			ratio1 := 100*est1/float64(expected[i]) - 100
			fmt.Printf("\n%s:\t\texpected %d\t\thlg ~= %.2f%%", id, expected[i], ratio1)
		}
	}
	fmt.Println(time.Since(now))
}
