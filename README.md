# FlipCounter

Data-structure with 5 byte keys and 3 bytes counters (64 bits per entry).

The 3 byte counter (24 bit) uses 23 for counting (guaranteeing 100% accurate counting per key up to `(2^23)-1 ==> 8388607` hits.

The remaining bit is used to annotate an overflow (hits > `8388607`) and the counter falls back to estimated counting trying to guarantee a max 1% absolute error.

## Usage is straight forward

```go
import "github.com/watchly/flipcounter"

c := flipcounter.New()
fc.Increment([]byte("watchly"))
fc.Get([]byte("watchly"))
```
