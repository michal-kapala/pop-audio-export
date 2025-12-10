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
