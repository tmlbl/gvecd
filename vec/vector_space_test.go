package vec

import (
	"fmt"
	"sort"
	"testing"
	"time"
)

func TestSearchVectorSpace(t *testing.T) {
	n := uint64(100_000)
	dim := uint64(300)
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

	// Create a new random point to search for
	needle := Vector{
		key:    "needle",
		values: randomValues(300),
	}
	start = time.Now()
	i := Search(vs, needle)
	fmt.Println("found nearest in", time.Since(start))
	resultDist := needle.distance(vs.slice[i])
	fmt.Println("distance between", needle.key, "and", vs.slice[i].key,
		"is", resultDist)

	// Confirm this is the nearest match
	for i, v := range vs.slice {
		dist := needle.distance(v)
		if abs(dist) < resultDist {
			t.Error(fmt.Errorf("result distance is %f but %d has distance of %f",
				resultDist, i, dist))
		}
	}
}
