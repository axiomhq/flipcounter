# FlipCounter

Data-structure with 5 byte keys and 3 bytes counters (64 bits per entry).
The 3 byte counter (24 bit) uses 23 for counting (guaranterrin 100% accurate counting per key up to `(2^23)-1 ==> 8388607` hits.
The remaining bit is used to signal estimated counting. It is flipped once counter > `8388607`, and estimation tries guaraneet a max 1% absolute error.