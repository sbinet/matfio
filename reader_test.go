package matfio

import (
	"io"
	"os"
	"testing"
)

func TestReader(t *testing.T) {
	for _, fname := range []string{
		"testdata/ex3weights.mat",
		"testdata/xy.mat",
		"testdata/big_endian.mat",
		"testdata/little_endian.mat",
		"testdata/teststruct_6.1_SOL2.mat",
	} {
		t.Logf("testing [%s]...\n", fname)
		f, err := os.Open(fname)
		if err != nil {
			t.Errorf("file: %s: %v\n", fname, err)
			continue
		}
		defer f.Close()

		r, err := NewReader(f)
		if err != nil {
			t.Errorf("file: %s: %v\n", fname, err)
			continue
		}
		if r != nil {
			println("r=", r)
		}

	dataLoop:
		for {
			var data DataElement
			err = r.Read(&data)
			if err == io.EOF {
				break dataLoop
			}
			if err != nil {
				t.Errorf("file: %s: %v\n", fname, err)
				continue
			}
			t.Logf("data.Tag:  %v\n", data.Tag)
			t.Logf("data.Size: %v\n", data.Size)
		}
	}
}
