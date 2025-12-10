package forge

import (
	"bytes"
	"fmt"
	"pop-audio-export/utils"
)

type NameTableEntry struct {
	DataSize          uint32
	FileDataId        uint64
	Unknown1          uint32
	ResourceId        uint32
	Unknown2          uint32 // Unknown2 + Unknown3 are a list of uint32s according to Blacksmith impl
	Unknown3          uint32
	NextFileCount     uint32
	PreviousFileCount uint32
	Unknown4          uint32
	Timestamp         uint32
	Name              string
	Unknown5          uint32 // Unknown5-8 are a list of uint32s according to Blacksmith impl (originally 5)
	Unknown6          uint32
	Unknown7          uint32
	Unknown8          uint32
}

func (entry *NameTableEntry) Read(reader *bytes.Reader) {
	size, err := utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("NameTableEntry.DataSize: %v", err))
	}
	entry.DataSize = size
	id, err := utils.ReadU64(reader)
	if err != nil {
		panic(fmt.Sprintf("NameTableEntry.FileDataId: %v", err))
	}
	entry.FileDataId = id
	unk32, err := utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("NameTableEntry.Unknown1: %v", err))
	}
	entry.Unknown1 = unk32
	resId, err := utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("NameTableEntry.ResourceId: %v", err))
	}
	entry.ResourceId = resId
	unk32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("NameTableEntry.Unknown2: %v", err))
	}
	entry.Unknown2 = unk32
	unk32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("NameTableEntry.Unknown3: %v", err))
	}
	entry.Unknown3 = unk32
	fileCount, err := utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("NameTableEntry.NextFileCount: %v", err))
	}
	entry.NextFileCount = fileCount
	fileCount, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("NameTableEntry.PreviousFileCount: %v", err))
	}
	entry.PreviousFileCount = fileCount
	unk32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("NameTableEntry.Unknown4: %v", err))
	}
	entry.Unknown4 = unk32
	timestamp, err := utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("NameTableEntry.Timestamp: %v", err))
	}
	entry.Timestamp = timestamp
	name, err := utils.ReadSizedString(reader, 128)
	if err != nil {
		panic(fmt.Sprintf("NameTableEntry.Name: %v", err))
	}
	entry.Name = name
	unk32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("NameTableEntry.Unknown5: %v", err))
	}
	entry.Unknown5 = unk32
	unk32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("NameTableEntry.Unknown6: %v", err))
	}
	entry.Unknown6 = unk32
	unk32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("NameTableEntry.Unknown7: %v", err))
	}
	entry.Unknown7 = unk32
	unk32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("NameTableEntry.Unknown8: %v", err))
	}
	entry.Unknown8 = unk32
}
