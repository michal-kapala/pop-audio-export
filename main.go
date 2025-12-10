package main

import (
	"fmt"
	"os"
	"path/filepath"
	"pop-audio-export/forge"
	"strings"

	"github.com/sqweek/dialog"
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
	// read forge
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	forgeFile, err := forge.Read(data)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	dir := strings.Replace(filepath.Base(path), filepath.Ext(path), "", 1)
	// create export directory
	if err = os.RemoveAll(dir); err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	err = os.Mkdir(dir, 0755)
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	// export files
	fmt.Printf("[FORGE]\tExporting %d files to %s\n", len(forgeFile.Data.Files), dir)
	for idx, file := range forgeFile.Data.Files {
		fmt.Printf("%d\t%s\n", idx+1, file.Header.FileName)
		file.Export(dir)
	}
}
