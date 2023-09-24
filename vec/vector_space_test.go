package vec

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestSearchVectorSpace(t *testing.T) {
	n := uint64(100000)
	dim := uint64(30)
	vs := randomVectorSpace(n, dim)

	start := time.Now()
	sort.Sort(vs)
	fmt.Println("sorted", n, "records in", time.Since(start))

	vs.Add(Vector{
		key:    "add another",
		values: randomValues(dim),
	})

	start = time.Now()
	sort.Sort(vs)
	fmt.Println("sorted new key in", time.Since(start))

	// Make sure they're in order
	root := vs.slice[0]
	var dist float32
	for i := range vs.slice {
		d := root.distance(vs.slice[i])
		if d > dist {
			// t.Errorf("wrong distance: %f > %f", d, dist)
		}
		dist = d
	}

	start = time.Now()
	needle := vs.slice[rand.Intn(len(vs.slice))]
	i := Search(vs, needle)
	fmt.Println("found nearest in", time.Since(start))
	fmt.Println("distance between", needle.key, "and", vs.slice[i].key,
		"is", needle.distance(vs.slice[i]))
}
