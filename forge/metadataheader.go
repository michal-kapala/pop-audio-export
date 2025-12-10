package forge

import (
	"bytes"
	"fmt"
	"pop-audio-export/utils"
)

type ForgeMetadataHeader struct {
	EntryCount       uint32
	Unknown1         uint32 // Unknown1 + Unknown2 are a list of uint32s according to Blacksmith impl; 1
	Unknown2         uint64 // 0
	Unknown3         uint64 // -1
	MaxFilesForIndex uint32 // 5k
	Unknown4         uint32 // 2
	DataOffset       uint64 // First data section offset
}

func (fmh *ForgeMetadataHeader) Read(reader *bytes.Reader) {
	count, err := utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("ForgeMetadataHeader.EntryCount: %v", err))
	}
	fmh.EntryCount = count
	unk32, err := utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("ForgeMetadataHeader.Unknown1: %v", err))
	}
	fmh.Unknown1 = unk32
	unk64, err := utils.ReadU64(reader)
	if err != nil {
		panic(fmt.Sprintf("ForgeMetadataHeader.Unknown2: %v", err))
	}
	fmh.Unknown2 = unk64
	unk64, err = utils.ReadU64(reader)
	if err != nil {
		panic(fmt.Sprintf("ForgeMetadataHeader.Unknown3: %v", err))
	}
	fmh.Unknown3 = unk64
	maxFiles, err := utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("ForgeMetadataHeader.MaxFilesForIndex: %v", err))
	}
	fmh.MaxFilesForIndex = maxFiles
	unk32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("ForgeMetadataHeader.Unknown4: %v", err))
	}
	fmh.Unknown4 = unk32
	off, err := utils.ReadU64(reader)
	if err != nil {
		panic(fmt.Sprintf("ForgeMetadataHeader.DataOffset: %v", err))
	}
	fmh.DataOffset = off
}
