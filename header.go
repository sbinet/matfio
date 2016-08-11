package matfio

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"unsafe"
)

type Tag uint32

//go:generate stringer -type Tag
const (
	Int8Tag       Tag = 1
	Uint8Tag      Tag = 2
	Int16Tag      Tag = 3
	Uint16Tag     Tag = 4
	Int32Tag      Tag = 5
	Uint32Tag     Tag = 6
	Float32Tag    Tag = 7
	_reservedTag1 Tag = 8
	Float64Tag    Tag = 9
	_reservedTag2 Tag = 10
	_reservedTag3 Tag = 11
	Int64Tag      Tag = 12
	Uint64Tag     Tag = 13
	MatrixTag     Tag = 14
	CompressedTag Tag = 15
	UTF8Tag       Tag = 16
	UTF16Tag      Tag = 17
	UTF32Tag      Tag = 18
)

type Class byte

//go:generate stringer -type Class
const (
	CellClass    Class = 1
	StructClass  Class = 2
	ObjectClass  Class = 3
	CharClass    Class = 4
	SparseClass  Class = 5
	Float64Class Class = 6
	Float32Class Class = 7
	Int8Class    Class = 8
	Uint8Class   Class = 9
	Int16Class   Class = 10
	Uint16Class  Class = 11
	Int32Class   Class = 12
	Uint32Class  Class = 13
	Int64Class   Class = 14
	Uint64Class  Class = 15
)

type Header struct {
	txt     [116]byte
	offset  int64
	version [2]byte
	endian  [2]byte
	order   binary.ByteOrder
}

func (hdr *Header) read(r io.ReaderAt) error {
	var pos int64
	n, err := r.ReadAt(hdr.txt[:], pos)
	if err != nil {
		return err
	}
	pos += int64(n)

	var subsys [8]byte
	n, err = r.ReadAt(subsys[:], pos)
	if err != nil {
		return err
	}
	pos += int64(n)

	n, err = r.ReadAt(hdr.version[:], pos)
	if err != nil {
		return err
	}
	pos += int64(n)

	n, err = r.ReadAt(hdr.endian[:], pos)
	if err != nil {
		return err
	}
	pos += int64(n)

	hdr.order = mlabOrder(hdr.endian)
	switch hdr.order {
	case binary.LittleEndian:
		hdr.version[0], hdr.version[1] = hdr.version[1], hdr.version[0]
	case binary.BigEndian:
		// no-op
	}

	fmt.Printf("txt:     %q\n", hdr.Text())
	fmt.Printf("subsys:  %q\n", subsys)
	fmt.Printf("version: %d (%x)\n", hdr.version, hdr.version)
	fmt.Printf("endian:  %q (%x) big-endian=%v\n", hdr.endian, hdr.endian, hdr.order)
	return err
}

func (hdr *Header) Text() string {
	txt := strings.TrimSpace(string(hdr.txt[:]))
	txt = strings.TrimRight(txt, "\x00")
	return txt
}

func mlabOrder(order [2]byte) binary.ByteOrder {
	const (
		MI uint16 = 0x4D49 // "MI"
		IM uint16 = 0x494D // "IM"
	)
	switch {
	case binary.LittleEndian.Uint16(order[:]) == MI:
		return binary.LittleEndian
	case binary.BigEndian.Uint16(order[:]) == MI:
		return binary.BigEndian
	}
	panic(fmt.Errorf("matfio: invalid endianness %[1]q (%[1]x)", order))
}

func isBigEndian() bool {
	val := int16(1)
	arr := *(*[2]byte)(unsafe.Pointer(&val))
	return binary.BigEndian.Uint16(arr[:]) == 1
}

type DataElement struct {
	Tag  Tag
	Size uint32
	Data []byte
}
