# Cuckoo
Implementing a Cuckoo filter, based on the paper Cuckoo Filter: Practically Better Than Bloom

![cuckoo](https://www.offthegridnews.com/wp-content/uploads/2015/08/rooster-fameimagesDOTcom.jpg)


## Example Usage

```go
n1 := []byte("Bess")
f := NewCuckooFilter(uint(b.N))
err := f.Insert(n1)
isPresnt := f.Lookup(n1)
```



## Performance

```bash
			BenchmarkCuckooAdd-4          	 2000000	      1008 ns/op
			BenchmarkCuckooTest-4         	 2000000	       685 ns/op
			BenchmarkCuckooTestAndAdd-4   	 1000000	      1780 ns/op
			BenchmarkCuckooLookupAndDelete-4 2000000	       796 ns/op
```


## Inspired
- https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf
- https://brilliant.org/wiki/cuckoo-filter/
- https://github.com/tylertreat/BoomFilters

## TODO

- Make it thread safe
- Improve hashing time, need to optimize on Insertion function
