package utils

func WriteU32(buf []byte, u32 uint32) []byte {
	slice := []byte{
		byte(u32),
		byte(u32 >> 8),
		byte(u32 >> 16),
		byte(u32 >> 24),
	}
	return append(buf, slice...)
}

func WriteU32BE(buf []byte, u32 uint32) []byte {
	slice := []byte{
		byte(u32 >> 24),
		byte(u32 >> 16),
		byte(u32 >> 8),
		byte(u32),
	}
	return append(buf, slice...)
}
