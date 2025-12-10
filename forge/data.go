package forge

import (
	"bytes"
	"fmt"
	"io"
)

type ForgeData struct {
	Header       ForgeMetadataHeader
	DataSections []*ForgeDataSection
	Files        []*FileData
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
	// files
	for _, section := range fd.DataSections {
		for _, entry := range section.IndexTable {
			_, err := reader.Seek(int64(entry.DataOffset), io.SeekStart)
			if err != nil {
				panic(fmt.Sprintf("IndexTableEntry.DataOffset: %v", err))
			}
			file := &FileData{}
			file.Read(reader)
			fd.Files = append(fd.Files, file)
		}
	}
}
