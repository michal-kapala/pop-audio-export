package utils

import (
	"bytes"
	"encoding/binary"
)

func ReadU32(reader *bytes.Reader) (uint32, error) {
	var result uint32
	err := binary.Read(reader, binary.LittleEndian, &result)
	return result, err
}

func ReadU64(reader *bytes.Reader) (uint64, error) {
	var result uint64
	err := binary.Read(reader, binary.LittleEndian, &result)
	return result, err
}

// ReadString is for null-terminated strings.
func ReadString(reader *bytes.Reader) (string, error) {
	var buf []byte
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return "", err
		}
		if b == 0 {
			break
		}
		buf = append(buf, b)
	}
	return string(buf), nil
}

// ReadSizedString is for statically-sized strings (char[x]).
func ReadSizedString(reader *bytes.Reader, size int) (string, error) {
	var buf []byte
	terminated := false
	for range size {
		b, err := reader.ReadByte()
		if err != nil {
			return "", err
		}
		if b == 0 {
			terminated = true
		}
		if !terminated {
			buf = append(buf, b)
		}
	}
	buf = append(buf, 0)
	return string(buf), nil
}
