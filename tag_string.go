// generated by stringer -type Tag; DO NOT EDIT

package matfio

import "fmt"

const _Tag_name = "Int8TagUint8TagInt16TagUint16TagInt32TagUint32TagFloat32Tag_reservedTag1Float64Tag_reservedTag2_reservedTag3Int64TagUint64TagMatrixTagCompressedTagUTF8TagUTF16TagUTF32Tag"

var _Tag_index = [...]uint8{0, 7, 15, 23, 32, 40, 49, 59, 72, 82, 95, 108, 116, 125, 134, 147, 154, 162, 170}

func (i Tag) String() string {
	i -= 1
	if i >= Tag(len(_Tag_index)-1) {
		return fmt.Sprintf("Tag(%d)", i+1)
	}
	return _Tag_name[_Tag_index[i]:_Tag_index[i+1]]
}
