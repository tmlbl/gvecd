package vec

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// DiskVectorSpace represents an immutable open file containing a sorted
// VectorSpace.
type DiskVectorSpace struct {
	r io.ReadSeeker
	Header
}

func (d *DiskVectorSpace) Len() int {
	return int(d.Header.len)
}

func (d *DiskVectorSpace) Get(i int) Vector {
	vec := Vector{
		values: make([]float32, d.dim),
	}
	// calculate the seek index to read from
	offset := d.startPos + int64(binary.Size(vec.values[0])*int(d.dim)*i)
	_, err := d.r.Seek(offset-1, io.SeekStart)
	if err != nil {
		panic(err)
	}

	for j := range vec.values {
		var val float32
		err = binary.Read(d.r, binary.NativeEndian, &val)
		if err != nil {
			panic(err)
		}
		vec.values[j] = val
		fmt.Println(val)
	}
	return vec
}

// OpenFile attempts to open a file as a DiskVectorSpace
func OpenFile(name string) (*DiskVectorSpace, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	header, err := readHeader(file)
	if err != nil {
		return nil, err
	}
	return &DiskVectorSpace{
		r:      file,
		Header: *header,
	}, nil
}
