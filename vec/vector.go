package vec

type Vector struct {
	key    string
	values []float32
}

func (v Vector) distance(q Vector) float32 {
	var sum float32
	for i, val := range v.values {
		d := val - q.values[i]
		sum += d
	}
	return sum
}
