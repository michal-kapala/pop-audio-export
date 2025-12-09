package bao

import (
	"bytes"
	"fmt"
	"pop-audio-export/utils"
)

// BAO (binary audio object, .bao/.sbao) file used by Ubisoft's proprietary sound engine Dare.
//
// For more information on the format refer to https://github.com/vgmstream/vgmstream/blob/ed976476635829ecb23b26b074a0c03ecabd0f7a/src/meta/ubi_bao.c
type BaoFile struct {
	Header *BaoHeader
	Data   *[]byte
	IsOgg  bool // PoP 2008 often simply stores raw .ogg files as data
}

func (bf *BaoFile) Read(reader *bytes.Reader, size int) {
	bf.Header = &BaoHeader{}
	bf.Header.Read(reader)
	data, err := utils.ReadBuffer(reader, size-int(bf.Header.HeaderSize))
	if err != nil {
		panic(fmt.Sprintf("BaoFile.Data: %v", err))
	}
	bf.Data = &data
	bf.IsOgg = string(data[:4]) == "OggS"
	if bf.IsOgg {
		fmt.Println("[BAO]\t\t^ is .ogg")
	}
}

type BaoHeader struct {
	Version    uint32 // BE, PoP 2008 = 0x021F0010
	HeaderSize uint32 // 40 bytes
	GUID       [16]byte
	Unknown1   uint32 // 0
	Unknown2   uint32 // 0
	Class      uint32 // 0x50000000
	Config     uint32 // 2
}

func (bh *BaoHeader) Read(reader *bytes.Reader) {
	u32, err := utils.ReadU32BE(reader)
	if err != nil {
		panic(fmt.Sprintf("BaoHeader.Version: %v", err))
	}
	bh.Version = u32
	u32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("BaoHeader.HeaderSize: %v", err))
	}
	bh.HeaderSize = u32
	guid, err := utils.ReadBuffer(reader, 16)
	if err != nil {
		panic(fmt.Sprintf("BaoHeader.GUID: %v", err))
	}
	bh.GUID = [16]byte(guid)
	u32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("BaoHeader.Unknown1: %v", err))
	}
	bh.Unknown1 = u32
	u32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("BaoHeader.Unknown2: %v", err))
	}
	bh.Unknown2 = u32
	u32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("BaoHeader.Class: %v", err))
	}
	bh.Class = u32
	u32, err = utils.ReadU32(reader)
	if err != nil {
		panic(fmt.Sprintf("BaoHeader.Config: %v", err))
	}
	bh.Config = u32
}
