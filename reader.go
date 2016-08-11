package matfio

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	hdrLvl5Size = 128
)

type Reader struct {
	r   io.ReaderAt
	hdr Header
	pos int64
	err error
}

func NewReader(r io.ReaderAt) (*Reader, error) {
	var hdr Header
	err := hdr.read(r)
	if err != nil {
		return nil, err
	}
	return &Reader{r: r, hdr: hdr, pos: hdrLvl5Size}, nil
}

func (r *Reader) Read(data *DataElement) error {
	if r.err != nil {
		return r.err
	}

	var tag [4]byte
	_, r.err = r.r.ReadAt(tag[:], r.pos)
	if r.err != nil {
		return r.err
	}
	fmt.Printf("tag[2]: %q\n", tag)
	// FIXME: handle small data element format
	tag = [4]byte{0, 0, 0, 0}
	switch tag {
	case [4]byte{0, 0, 0, 0}:
		data.Tag = Tag(r.readU32())
		data.Size = r.readU32()
		sz := data.Size
		// FIXME(sbinet): handle padding bytes
		if data.Tag != MatrixTag && false {
			sz = alignU32(sz)
		}
		data.Data = make([]byte, sz)

		n, err := r.r.ReadAt(data.Data, r.pos)
		if err != nil {
			r.err = err
		}
		r.pos += int64(n)
		data.Data = data.Data[:data.Size]
		fmt.Printf("tag: %v\n", data.Tag)
		fmt.Printf("siz: %v (%v)\n", data.Size, alignU32(data.Size))
	default:
		// small data element format
	}

	if data.Tag == CompressedTag {
		var rz io.ReadCloser
		rz, r.err = zlib.NewReader(bytes.NewReader(data.Data))
		if r.err != nil {
			return r.err
		}
		defer rz.Close()

		de, err := readDE(rz, r.hdr.order)
		if err != nil {
			r.err = err
			return r.err
		}
		*data = de
		fmt.Printf("=> tag: %v\n", data.Tag)
		fmt.Printf("=> siz: %v\n", data.Size)
	}

	if data.Tag == MatrixTag {
		rr := bytes.NewReader(data.Data)
		var err error
		var de DataElement
	matrixLoop:
		for {
			de, err = readDE(rr, r.hdr.order)
			if err != nil {
				break matrixLoop
			}
			fmt.Printf("===> tag: %v\n", de.Tag)
			fmt.Printf("===> siz: %v\n", de.Size)

		}
		if err == io.EOF {
			err = nil
		}
		if err != nil {
			r.err = err
			return r.err
		}
	}

	return r.err
}

func (r *Reader) readU32() (val uint32) {
	if r.err != nil {
		return
	}
	var buf [4]byte
	var n int
	n, r.err = r.r.ReadAt(buf[:], r.pos)
	if r.err != nil {
		return
	}
	r.pos += int64(n)

	val = r.hdr.order.Uint32(buf[:])
	return
}

func readDE(r io.Reader, bo binary.ByteOrder) (DataElement, error) {
	var (
		err error
		de  DataElement
		tag [4]byte
	)

	// FIXME(sbinet): handle small data element format
	switch tag {
	case [4]byte{0, 0, 0, 0}:
		var v uint32
		v, err = readU32(r, bo)
		if err != nil {
			return de, err
		}
		de.Tag = Tag(v)

		de.Size, err = readU32(r, bo)
		if err != nil {
			return de, err
		}
		de.Data = make([]byte, de.Size)
		_, err = r.Read(de.Data)
		if err != nil {
			return de, err
		}
	default:
		// small data element format
	}

	return de, err
}

func readU32(r io.Reader, bo binary.ByteOrder) (uint32, error) {
	var buf [4]byte
	_, err := r.Read(buf[:])
	if err != nil {
		return 0, err
	}
	return bo.Uint32(buf[:]), nil
}
