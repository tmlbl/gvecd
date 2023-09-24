package vec

import (
	"bytes"
	"testing"
)

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
