package vec

import (
	"fmt"
	"math/rand"
)

func randomValues(dim uint64) []float32 {
	values := make([]float32, dim)
	for j := uint64(0); j < dim; j++ {
		values[j] = rand.Float32()
	}
	return values
}

func randomVectorSpace(n, dim uint64) *MemVectorSpace {
	vs := NewMemVectorSpace(dim)
	for i := uint64(0); i < n; i++ {
		vs.Add(Vector{
			key:    fmt.Sprintf("key-%d", i),
			values: randomValues(dim),
		})
	}
	return vs
}
