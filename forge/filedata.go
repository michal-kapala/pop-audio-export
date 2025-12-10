package forge

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"pop-audio-export/bao"
	"pop-audio-export/utils"
	"strings"
)

// FILEDATA entry.
type FileData struct {
	Header FileDataHeader
	Data   *[]byte // Non-BAO file
	Bao    *bao.BaoFile
}

func (fd *FileData) Read(reader *bytes.Reader) {
	fd.Header.Read(reader)
	fd.Bao = &bao.BaoFile{}
	// name check for BAO since Header.ClassID is always 0 in PoP 2008
	if strings.Contains(fd.Header.FileName, "_BAO_") {
		fd.Bao.Read(reader, int(fd.Header.FileLength))
		fd.Data = nil
	} else {
		fd.Bao = nil
		data, err := utils.ReadBuffer(reader, int(fd.Header.FileLength))
		if err != nil {
			panic(fmt.Sprintf("FileData.Data: %v", err))
		}
		fd.Data = &data
	}
}

func (fd *FileData) Export(dir string) {
	// BAO file
	if fd.Bao != nil {
		data, ext := fd.Bao.Export()
		file, err := os.Create(filepath.Join(dir, fmt.Sprintf("%s%s", fd.Header.FileName, ext)))
		if err != nil {
			panic(fmt.Sprintf("FileData.Export: %v", err))
		}
		defer file.Close()
		_, err = file.Write(data)
		if err != nil {
			panic(fmt.Sprintf("FileData.Export: %v", err))
		}
		return
	}

	if fd.Data == nil {
		panic("FileData.Export: no available buffer for export")
	}
	// non-BAO file
	file, err := os.Create(filepath.Join(dir, fmt.Sprintf("%s.data", fd.Header.FileName)))
	if err != nil {
		panic(fmt.Sprintf("FileData.Export: %v", err))
	}
	defer file.Close()
	_, err = file.Write(*fd.Data)
	if err != nil {
		panic(fmt.Sprintf("FileData.Export: %v", err))
	}
}

// Originally scimitar::FatFileDataHeader.
type FileDataHeader struct {
	Magic                    string // FILEDATA
	FileName                 string // char[128]
	FilePath                 string // char[255]
	FileKey                  uint32 // scimitar::BigFileKey; referenced in index table/file name
	FileLength               uint32 // Total size of data excluding this header
	UMACHash                 uint64 // Referenced in the name table
	RevisionNumberData       uint32 // 0
	RevisionNumberAttributes uint32 // 0
	FileFlags                uint32 // 1 in GlobalMetaFile, 0 in asset files
	SCCStatusData            uint32 // scimitar::BigFileSCCAction; 2 in GlobalMetaFile, 4 in asset files
	SCCStatusAttributes      uint32 // scimitar::BigFileSCCAction; 2 in GlobalMetaFile, 4 in asset files
	AssociatedMetaFileKey    uint32 // scimitar::BigFileKey; 0
	Time                     uint32 // The same as in name table
	Deleted                  bool
	ClassID                  uint32
}

func (fdh *FileDataHeader) Read(reader *bytes.Reader) {
	str, err := utils.ReadSizedString(reader, 8)
	if err != nil {
		panic(fmt.Sprintf("FileDataHeader.Magic: %v", err))
	}
	if str != "FILEDATA" {
		panic(fmt.Sprintf("FileDataHeader.Magic: invalid FILEDATA entry signature: %s", str))
	}
	fdh.Magic = str
	str, err = utils.ReadSizedString(reader, 128)
	if err != nil {
		panic(fmt.Sprintf("FileDataHeader.FileName: %v", err))
	}
	fdh.FileName = str
	str, err = utils.ReadSizedString(reader, 255)
	if err != nil {
		panic(fmt.Sprintf("FileDataHeader.FilePath: %v", err))
	}
	fdh.FilePath = str
	u32, err := utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("FileDataHeader.FileKey: %v", err))
	}
	fdh.FileKey = u32
	len, err := utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("FileDataHeader.FileLength: %v", err))
	}
	fdh.FileLength = len
	hash, err := utils.ReadU64(reader)
	if err != nil {
		panic(fmt.Sprintf("FileDataHeader.UMACHash: %v", err))
	}
	fdh.UMACHash = hash
	u32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("FileDataHeader.RevisionNumberData: %v", err))
	}
	fdh.RevisionNumberData = u32
	u32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("FileDataHeader.RevisionNumberAttributes: %v", err))
	}
	fdh.RevisionNumberAttributes = u32
	u32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("FileDataHeader.FileFlags: %v", err))
	}
	fdh.FileFlags = u32
	u32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("FileDataHeader.SCCStatusData: %v", err))
	}
	fdh.SCCStatusData = u32
	u32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("FileDataHeader.SCCStatusAttributes: %v", err))
	}
	fdh.SCCStatusAttributes = u32
	u32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("FileDataHeader.AssociatedMetaFileKey: %v", err))
	}
	fdh.AssociatedMetaFileKey = u32
	u32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("FileDataHeader.Time: %v", err))
	}
	fdh.Time = u32
	deleted, err := reader.ReadByte()
	if err != nil {
		panic(fmt.Sprintf("FileDataHeader.Deleted: %v", err))
	}
	fdh.Deleted = deleted != 0
	u32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("FileDataHeader.ClassID: %v", err))
	}
	fdh.ClassID = u32
}
