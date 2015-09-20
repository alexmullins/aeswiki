AESWiki scrapes the AES Sbox Wikipedia page checking the
SBox/Inverse table's bytes match the given bytes on the page.
https://en.wikipedia.org/wiki/Rijndael_S-box

Uses goquery to scrape the <pre> tags. 

Possible speed ups:
1. Set initial capacity for the buffers returned in parseSBox* functions (done)
2. Use bytes from (Scanner).Bytes() to reduce allocations instead of using strings
3. Remove the prees map and run each parseSBox* in parallel

Profiling 
1. aeswiki.bench0 - Benchmark stats after completing (#1) above
PASS
BenchmarkParseSBoxTable-4	  100000	     21976 ns/op	    7200 B/op	      83 allocs/op
BenchmarkParseSBoxBytes-4	   30000	     53323 ns/op	   12080 B/op	     132 allocs/op
ok  	github.com/alexmullins/aeswiki	4.577s

