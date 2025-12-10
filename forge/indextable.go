package forge

import (
	"bytes"
	"fmt"
	"pop-audio-export/utils"
)

type IndexTableEntry struct {
	DataOffset uint64
	FileKey    uint32
	DataSize   uint32
}

func (entry *IndexTableEntry) Read(reader *bytes.Reader) {
	off, err := utils.ReadU64(reader)
	if err != nil {
		panic(fmt.Sprintf("IndexTableEntry.DataOffset: %v", err))
	}
	entry.DataOffset = off
	id, err := utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("IndexTableEntry.FileDataId: %v", err))
	}
	entry.FileKey = id
	size, err := utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("IndexTableEntry.DataSize: %v", err))
	}
	entry.DataSize = size
}
