package vec

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type MemVectorSpace struct {
	slice []Vector
	dim   uint64
}

func NewMemVectorSpace(dim uint64) *MemVectorSpace {
	return &MemVectorSpace{
		dim: dim,
	}
}

// Add validates the values and appends them to the slice
func (vs *MemVectorSpace) Add(v Vector) error {
	if uint64(len(v.values)) != vs.dim {
		return errors.New("vector of wrong length")
	}
	for _, x := range v.values {
		if x < 0 || x > 1 {
			return fmt.Errorf("values must be between 0 and 1")
		}
	}
	vs.slice = append(vs.slice, v)
	return nil
}

func (vs *MemVectorSpace) Len() int {
	return len(vs.slice)
}

func (vs *MemVectorSpace) Less(i, j int) bool {
	dist := vs.slice[i].distance(vs.slice[j])
	return dist < 0
}

func (vs *MemVectorSpace) Swap(i, j int) {
	iv := vs.slice[i]
	jv := vs.slice[j]
	vs.slice[i] = jv
	vs.slice[j] = iv
}

func (vs *MemVectorSpace) Get(i int) Vector {
	return vs.slice[i]
}

// Write serializes the VectorSpace to a searchable binary format
func (vs *MemVectorSpace) Write(w io.Writer) error {
	// Get the total length of the metadata and keys to determine the seek start
	// position for reading vector data
	var seekStart uint64 = binary.MaxVarintLen64 * 3
	for _, v := range vs.slice {
		seekStart += uint64(len(v.key))
	}

	// Metadata section
	// Write the length and dimension of the VectorSpace
	err := binary.Write(w, binary.NativeEndian, seekStart)
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
	// Write vector data
	for _, vec := range vs.slice {
		for _, val := range vec.values {
			err = binary.Write(w, binary.NativeEndian, val)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Read deserializes a VectorSpace
func Read(r io.Reader) (*MemVectorSpace, error) {
	// Read metadata
	header, err := readHeader(r)
	if err != nil {
		return nil, err
	}

	vs := &MemVectorSpace{
		dim:   header.dim,
		slice: make([]Vector, int(header.len)),
	}

	// Read key data
	for i := range vs.slice {
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

		// Initialize values
		vs.slice[i].values = make([]float32, header.dim)
	}

	// Read vector data
	for i := range vs.slice {
		for j := range vs.slice[i].values {
			var val float32
			err = binary.Read(r, binary.NativeEndian, &val)
			if err != nil {
				return nil, err
			}
			vs.slice[i].values[j] = val
		}
	}

	return vs, nil
}
