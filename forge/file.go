package forge

import (
	"bytes"
	"fmt"
	"io"
	"pop-audio-export/utils"
)

type Readable interface {
	Read(reader *bytes.Reader)
}

type ForgeFile struct {
	Header ForgeHeader
	Data   ForgeData
}

type ForgeHeader struct {
	Signature             string
	FileVersionIdentifier uint32
	DataHeaderOffset      uint64
}

func (fh *ForgeHeader) Read(reader *bytes.Reader) {
	sig, err := utils.ReadString(reader)
	if err != nil {
		panic(fmt.Sprintf("ForgeHeader.Signature: %v", err))
	}
	if sig != "scimitar" {
		panic(fmt.Sprintf("ForgeHeader.Signature: not a forge signature: '%s', expected 'scimitar'", sig))
	}
	fh.Signature = sig
	versionId, err := utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("ForgeHeader.FileVersionIdentifier: %v", err))
	}
	fh.FileVersionIdentifier = versionId
	off, err := utils.ReadU64(reader)
	if err != nil {
		panic(fmt.Sprintf("ForgeHeader.DataHeaderOffset: %v", err))
	}
	fh.DataHeaderOffset = off
}

type ForgeData struct {
	Header       ForgeMetadataHeader
	DataSections []*ForgeDataSection
}

func (fd *ForgeData) Read(reader *bytes.Reader) {
	// data header
	fd.Header.Read(reader)
	// 1st data section
	_, err := reader.Seek(int64(fd.Header.DataOffset), io.SeekStart)
	if err != nil {
		panic(fmt.Sprintf("ForgeData.Header.DataOffset: %v", err))
	}
	ds := &ForgeDataSection{}
	ds.Read(reader)
	fd.DataSections = append(fd.DataSections, ds)
	// subsequent sections
	if int64(ds.NextSectionOffset) != -1 {
		_, err = reader.Seek(int64(ds.NextSectionOffset), io.SeekStart)
		if err != nil {
			panic(fmt.Sprintf("ForgeDataSection.NextSectionOffset: %v", err))
		}
		for {
			section := &ForgeDataSection{}
			section.Read(reader)
			fd.DataSections = append(fd.DataSections, section)
			if int64(section.NextSectionOffset) == -1 {
				break
			} else {
				_, err = reader.Seek(int64(ds.NextSectionOffset), io.SeekStart)
				if err != nil {
					panic(fmt.Sprintf("ForgeDataSection.NextSectionOffset: %v", err))
				}
			}
		}
	}
}

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

type IndexTableEntry struct {
	DataOffset uint64
	FileDataId uint32
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
	entry.FileDataId = id
	size, err := utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("IndexTableEntry.DataSize: %v", err))
	}
	entry.DataSize = size
}

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

type FileData struct {
	IndexEntry IndexTableEntry
	NameEntry  NameTableEntry
	Data       []byte
}

// Reads a .forge file into its structure.
func Read(file []byte) (*ForgeFile, error) {
	forge := &ForgeFile{}
	reader := bytes.NewReader(file)
	forge.Header.Read(reader)
	_, err := reader.Seek(int64(forge.Header.DataHeaderOffset), io.SeekStart)
	if err != nil {
		panic(fmt.Sprintf("ForgeFile.Header.DataHeaderOffset: %v", err))
	}
	forge.Data.Read(reader)
	return forge, nil
}
