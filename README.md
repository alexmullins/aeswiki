AESWiki scrapes the AES Sbox Wikipedia page checking the
SBox/Inverse tables bytes match the given C++ init bytes on the page.
https://en.wikipedia.org/wiki/Rijndael_S-box

Uses goquery to scrape the pre tags. 

Possible speed ups:
1. Set initial capacity for the buffers returned in parseSBox* functions (done)
2. Move away from strings and use bytes from (Scanner).Bytes() to reduce allocations (not needed)
3. Remove the prees map and run each parseSBox* in parallel 

Profiling 
1.aeswiki.bench0 - Benchmark stats after completing (#1) above

$ go test -bench=BenchmarkParse
PASS
BenchmarkParseSBoxTable-4	  100000	     21976 ns/op	    7200 B/op	      83 allocs/op
BenchmarkParseSBoxBytes-4	   30000	     53323 ns/op	   12080 B/op	     132 allocs/op
ok  	github.com/alexmullins/aeswiki	4.577s

2.After profiling both functions, most of the cumulative time was spent in the Replace func.
After testing the Replace() in the strings and bytes packages there was little to
no performance difference to warrant a change. 

$ go test -bench=BenchmarkReplace
PASS
BenchmarkReplaceString-4	10000000	       204 ns/op	      32 B/op	       2 allocs/op
BenchmarkReplaceByte-4  	10000000	       203 ns/op	      16 B/op	       1 allocs/op
ok  	github.com/alexmullins/aeswiki	4.513s

Above there is no timing difference, half the bytes/op and allocs/op. The string.Replace
is fine for this use case. 