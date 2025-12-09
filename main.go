package main

import (
	"os"
	//"path/filepath"
	"fmt"
	"github.com/sqweek/dialog"
	"pop-audio-export/forge"
)

func main() {
	fmt.Println("Choose .forge file:")
	path, err := dialog.File().Filter("Forge assets (.forge)", "forge").Load()
	if err != nil {
		if err == dialog.Cancelled {
			fmt.Println(err)
			return
		} else {
			panic(fmt.Sprintf("%v", err))
		}
	}
	fmt.Println(path)
	//fmt.Println(filepath.Base(path))
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	forgeFile, err := forge.Read(data)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	fmt.Printf("DataSections[0].IndexTable[0].DataSize: %x\n", forgeFile.Data.DataSections[0].IndexTable[0].DataSize)
	fmt.Printf("DataSections[1].IndexTable[1].DataSize: %x\n", forgeFile.Data.DataSections[1].IndexTable[0].DataSize)
	fmt.Printf("DataSections[0].IndexTable.Count: %d\n", len(forgeFile.Data.DataSections[0].IndexTable))
	fmt.Printf("DataSections[1].IndexTable.Count: %d\n", len(forgeFile.Data.DataSections[1].IndexTable))
	fmt.Printf("DataSections[0].NameTable[0].Name: %s\n", forgeFile.Data.DataSections[0].NameTable[0].Name)
	fmt.Printf("DataSections[1].NameTable[0].Name: %s\n", forgeFile.Data.DataSections[1].NameTable[0].Name)
}
