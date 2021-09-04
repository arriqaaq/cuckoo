# Cuckoo
Implementing a Cuckoo filter, based on the paper Cuckoo Filter: Practically Better Than Bloom



## Example Usage

```go
package main

import(
	"log"
	"github.com/arriqaaq/cuckoo"
)

func main(){
	f := cuckoo.NewCuckooFilter(uint(5000000), 0.001)
	n1 := []byte("Bess")
	err := f.Insert(n1)
	if err!=nil{
		log.Println("bucket full: ",err)
	}
	isPresent := f.Lookup(n1)
	log.Println("key present? ",isPresent)

}
```



## Performance

```bash
			BenchmarkCuckooAdd-4          	 2000000	      1008 ns/op
			BenchmarkCuckooTest-4         	 2000000	       685 ns/op
			BenchmarkCuckooTestAndAdd-4   	 1000000	      1780 ns/op
			BenchmarkCuckooLookupAndDelete-4 2000000	       796 ns/op
```


## Reference
- https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf
- https://brilliant.org/wiki/cuckoo-filter/
- https://github.com/tylertreat/BoomFilters

## TODO

- Make it thread safe
- Improve hashing time, need to optimize on Insertion function
