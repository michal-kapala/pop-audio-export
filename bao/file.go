package bao

import (
	"bytes"
	"fmt"
	"pop-audio-export/utils"
)

// BAO (binary audio object, .bao/.sbao) file used by Ubisoft's proprietary sound engine Dare.
//
// For more information on the format, refer to https://github.com/vgmstream/vgmstream/blob/ed976476635829ecb23b26b074a0c03ecabd0f7a/src/meta/ubi_bao.c
type BaoFile struct {
	Header *BaoHeader
	Data   []byte
	IsOgg  bool // PoP 2008 often simply stores raw .ogg files as BAO data
}

func (bf *BaoFile) Read(reader *bytes.Reader, size int) {
	bf.Header = &BaoHeader{}
	bf.Header.Read(reader)
	data, err := utils.ReadBuffer(reader, size-int(bf.Header.HeaderSize))
	if err != nil {
		panic(fmt.Sprintf("BaoFile.Data: %v", err))
	}
	bf.Data = data
	bf.IsOgg = string(data[:4]) == "OggS"
}

// Returns a data buffer and extension.
func (bf *BaoFile) Export() ([]byte, string) {
	var buf []byte
	var extension string
	if bf.IsOgg {
		buf = bf.Data
		extension = ".ogg"
	} else {
		buf = bf.Header.Export()
		buf = append(buf, bf.Data...)
		extension = ".bao"
	}
	return buf, extension
}
