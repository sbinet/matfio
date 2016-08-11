package matfio

func alignU32(offset uint32) uint32 {
	const align = 8
	// return offset + (align-(offset%align))%align
	return offset + (align - (offset&(align-1))&(align-1))
}
