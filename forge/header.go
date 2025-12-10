package forge

import (
	"bytes"
	"fmt"
	"pop-audio-export/utils"
)

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
