package vec

// VectorSpace is the interface unifying in-memory and on-disk representations
// of vector spaces
type VectorSpace interface {
	Len() int
	Get(i int) Vector
}

// Search is a simple binary search on the sorted VectorSpace in memory.
// It returns the (usually exact) index of the nearest point.
func Search(vs VectorSpace, v Vector) int {
	i, div := vs.Len()/2, 2
	var dist float32
	for div <= vs.Len() {
		dist = v.distance(vs.Get(i))
		div *= 2
		if dist > 0 {
			i += vs.Len() / div
		} else {
			i -= vs.Len() / div
		}
	}
	if i > 0 && v.distance(vs.Get(i-1)) < dist {
		return i - 1
	} else if i < vs.Len()-1 && v.distance(vs.Get(i+1)) < dist {
		return i + 1
	}
	return i
}
