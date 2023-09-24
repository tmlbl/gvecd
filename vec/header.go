package vec

import (
	"encoding/binary"
	"io"
)

// Header contains metadata about a VectorSpace's on-disk representation
type Header struct {
	dim uint64
	len uint64
	// the seek index representing the start of the vector data section
	startPos int64
}

func readHeader(r io.Reader) (*Header, error) {
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

	return &Header{
		dim:      dim,
		len:      length,
		startPos: int64(mdSize),
	}, nil
}
