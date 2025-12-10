package bao

import (
	"bytes"
	"fmt"
	"pop-audio-export/utils"
)

// BAO header version for PoP 2008.
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

func (bh *BaoHeader) Export() []byte {
	var header []byte
	header = utils.WriteU32BE(header, bh.Version)
	header = utils.WriteU32(header, bh.HeaderSize)
	header = append(header, bh.GUID[:]...)
	header = utils.WriteU32(header, bh.Unknown1)
	header = utils.WriteU32(header, bh.Unknown2)
	header = utils.WriteU32(header, bh.Class)
	header = utils.WriteU32(header, bh.Config)
	return header
}
