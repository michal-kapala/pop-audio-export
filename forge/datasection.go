package forge

import (
	"bytes"
	"fmt"
	"io"
	"pop-audio-export/utils"
)

type ForgeDataSection struct {
	IndexCount            uint32 // Number of indexed entries in this section
	FollowingSectionCount uint32 // Excludes self, 0 if last
	IndexTableOffset      uint64 // Index table for this section
	NextSectionOffset     uint64 // Next data section, -1 if last
	IndexStart            uint32 // First entry index in this section
	IndexEnd              uint32 // Last possible index in this section
	NameTableOffset       uint64 // Name table for this section
	UnknownOffset         uint64 // After name table, a short metadata segment

	IndexTable []*IndexTableEntry
	NameTable  []*NameTableEntry
}

func (fds *ForgeDataSection) Read(reader *bytes.Reader) {
	idx, err := utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("ForgeDataSection.IndexCount: %v", err))
	}
	fds.IndexCount = idx
	sectionCount, err := utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("ForgeDataSection.FollowingSectionCount: %v", err))
	}
	fds.FollowingSectionCount = sectionCount
	off, err := utils.ReadU64(reader)
	if err != nil {
		panic(fmt.Sprintf("ForgeDataSection.IndexTableOffset: %v", err))
	}
	fds.IndexTableOffset = off
	off, err = utils.ReadU64(reader)
	if err != nil {
		panic(fmt.Sprintf("ForgeDataSection.NextSectionOffset: %v", err))
	}
	fds.NextSectionOffset = off
	idx, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("ForgeDataSection.IndexStart: %v", err))
	}
	fds.IndexStart = idx
	idx, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("ForgeDataSection.IndexEnd: %v", err))
	}
	fds.IndexEnd = idx
	off, err = utils.ReadU64(reader)
	if err != nil {
		panic(fmt.Sprintf("ForgeDataSection.NameTableOffset: %v", err))
	}
	fds.NameTableOffset = off
	off, err = utils.ReadU64(reader)
	if err != nil {
		panic(fmt.Sprintf("ForgeDataSection.UnknownOffset: %v", err))
	}
	fds.UnknownOffset = off
	// current stream position
	pos := reader.Size() - int64(reader.Len())
	// index table
	_, err = reader.Seek(int64(fds.IndexTableOffset), io.SeekStart)
	if err != nil {
		panic(fmt.Sprintf("ForgeDataSection.IndexTableOffset: %v", err))
	}
	for idx = fds.IndexStart; idx < fds.IndexStart+fds.IndexCount; idx++ {
		idxEntry := &IndexTableEntry{}
		idxEntry.Read(reader)
		fds.IndexTable = append(fds.IndexTable, idxEntry)
	}
	// name table
	_, err = reader.Seek(int64(fds.NameTableOffset), io.SeekStart)
	if err != nil {
		panic(fmt.Sprintf("ForgeDataSection.NameTableOffset: %v", err))
	}
	for idx = fds.IndexStart; idx < fds.IndexStart+fds.IndexCount; idx++ {
		nameEntry := &NameTableEntry{}
		nameEntry.Read(reader)
		fds.NameTable = append(fds.NameTable, nameEntry)
	}

	// restore stream position
	_, err = reader.Seek(pos, io.SeekStart)
	if err != nil {
		panic(fmt.Sprintf("ForgeDataSection: failed to restore stream: %v", err))
	}
}
