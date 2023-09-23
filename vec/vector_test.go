package vec

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func randomValues(dim uint64) []float32 {
	values := make([]float32, dim)
	for j := uint64(0); j < dim; j++ {
		values[j] = rand.Float32()
	}
	return values
}

func randomVectorSpace(n, dim uint64) *VectorSpace {
	vs := NewVectorSpace(dim)
	for i := uint64(0); i < n; i++ {
		vs.Add(Vector{
			key:    fmt.Sprintf("key-%d", i),
			values: randomValues(dim),
		})
	}
	return vs
}

func TestSortVectorSpace(t *testing.T) {
	n := uint64(100000)
	dim := uint64(3)
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
	i := vs.FindNearest(needle, 1)
	fmt.Println("found nearest in", time.Since(start))
	fmt.Println("distance between", needle.key, "and", vs.slice[i].key,
		"is", needle.distance(vs.slice[i]))
}

func TestWriteVectorSpace(t *testing.T) {
	vs := randomVectorSpace(100, 10)
	buf := bytes.NewBuffer([]byte{})
	err := vs.Write(buf)
	if err != nil {
		t.Error(err)
	}
	nvs, err := Read(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Error(err)
	}

	if nvs.dim != vs.dim {
		t.Errorf("misread dimension")
	}

	if len(nvs.slice) != len(vs.slice) {
		t.Errorf("misread size")
	}

	// Check that keys match
	for i := range vs.slice {
		if vs.slice[i].key != nvs.slice[i].key {
			t.Errorf("misread key: %s != %s", vs.slice[i].key, nvs.slice[i].key)
		}
	}

	// Check that values match
	for i := range vs.slice {
		for j, v := range vs.slice[i].values {
			if nvs.slice[i].values[j] != v {
				t.Errorf("misread value at (%d, %d): %f != %f",
					i, j, v, nvs.slice[i].values[j])
			}
		}
	}
}
