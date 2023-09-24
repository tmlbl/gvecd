package vec

import (
	"os"
	"testing"
)

func TestSearchDiskVectorSpace(t *testing.T) {
	vs := randomVectorSpace(10, 2)
	f, err := os.CreateTemp("", "gvec-test-")
	if err != nil {
		t.Error(err)
	}

	err = vs.Write(f)
	if err != nil {
		t.Error(err)
	}
	err = f.Close()
	if err != nil {
		t.Error(err)
	}

	dvs, err := OpenFile(f.Name())
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < vs.Len(); i++ {
		target := vs.Get(i)
		result := dvs.Get(i)
		for j := range target.values {
			if target.values[j] != result.values[j] {
				t.Errorf("mismatched values")
			}
		}
	}
}
