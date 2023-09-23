package vec

import (
	"encoding/binary"
	"errors"
	"io"
)

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

type VectorSpace struct {
	slice []Vector
	dim   uint64
}

func NewVectorSpace(dim uint64) *VectorSpace {
	return &VectorSpace{
		dim: dim,
	}
}

func (vs *VectorSpace) Add(v Vector) error {
	if uint64(len(v.values)) != vs.dim {
		return errors.New("vector of wrong length")
	}
	vs.slice = append(vs.slice, v)
	return nil
}

func (vs *VectorSpace) Len() int {
	return len(vs.slice)
}

func (vs *VectorSpace) Less(i, j int) bool {
	dist := vs.slice[i].distance(vs.slice[j])
	return dist < 0
}

func (vs *VectorSpace) Swap(i, j int) {
	iv := vs.slice[i]
	jv := vs.slice[j]
	vs.slice[i] = jv
	vs.slice[j] = iv
}

// FindNearest is a simple binary search on the sorted VectorSpace in memory
func (vs *VectorSpace) FindNearest(v Vector, results int) int {
	i, div := len(vs.slice)/2, 2
	var dist float32
	for div <= len(vs.slice) {
		dist = v.distance(vs.slice[i])
		div *= 2
		if dist > 0 {
			i += len(vs.slice) / div
		} else {
			i -= len(vs.slice) / div
		}
		// fmt.Println("distance at", i, "is", dist, "div is", div)
	}
	if i > 0 && v.distance(vs.slice[i-1]) < dist {
		return i - 1
	} else if i < len(vs.slice)-1 && v.distance(vs.slice[i+1]) < dist {
		return i + 1
	}
	return i
}

// Write serializes the VectorSpace to a searchable binary format
func (vs *VectorSpace) Write(w io.Writer) error {
	// Get all of the offsets for the key data so we know the total size
	// of the metadata/key section
	offsets := make([]uint64, len(vs.slice))
	var cur uint64
	for i, v := range vs.slice {
		offsets[i] = cur
		cur += uint64(len(v.key))
	}

	// Metadata section
	// Write the length and dimension of the VectorSpace
	err := binary.Write(w, binary.NativeEndian, cur)
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.NativeEndian, uint64(len(vs.slice)))
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.NativeEndian, uint64(vs.dim))
	if err != nil {
		return err
	}
	// Write the key data
	for _, v := range vs.slice {
		err = binary.Write(w, binary.NativeEndian, uint64(len(v.key)))
		if err != nil {
			return err
		}
		_, err := w.Write([]byte(v.key))
		if err != nil {
			return err
		}
	}
	// // Write vector data with offsets
	// for i, vec := range vs.slice {
	// 	for _, val := range vec.values {
	// 		err = binary.Write(w, binary.NativeEndian, val)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// 	err = binary.Write(w, binary.NativeEndian, offsets[i])
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

// Read deserializes a VectorSpace
func Read(r io.Reader) (*VectorSpace, error) {
	// Read metadata
	var mdSize uint64
	err := binary.Read(r, binary.NativeEndian, &mdSize)
	if err != nil {
		return nil, err
	}
	var length uint64
	err = binary.Read(r, binary.NativeEndian, &length)
	if err != nil {
		return nil, err
	}
	var dim uint64
	err = binary.Read(r, binary.NativeEndian, &dim)
	if err != nil {
		return nil, err
	}

	vs := &VectorSpace{
		dim:   dim,
		slice: make([]Vector, int(length)),
	}

	// Read key data
	i := 0
	for i < vs.Len() {
		var keyLen uint64
		err = binary.Read(r, binary.NativeEndian, &keyLen)
		if err != nil {
			return nil, err
		}
		key := make([]byte, keyLen)
		_, err := r.Read(key)
		if err != nil {
			return nil, err
		}
		vs.slice[i].key = string(key)
		i++
	}
	// pos := binary.MaxVarintLen64 * 3
	// for pos < int(mdSize) {
	// 	var keyLen uint64
	// 	err := binary.Read(r, binary.NativeEndian, &keyLen)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	fmt.Println("keyLen", keyLen)
	// 	pos += binary.MaxVarintLen64
	// 	key := make([]byte, keyLen)
	// 	n, err := r.Read(key)
	// 	pos += n
	// 	fmt.Println(string(key))
	// }

	return vs, nil
}
