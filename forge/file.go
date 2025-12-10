package forge

import (
	"bytes"
	"fmt"
	"io"
)

type ForgeFile struct {
	Header ForgeHeader
	Data   ForgeData
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
