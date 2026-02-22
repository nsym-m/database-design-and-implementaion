package main

import (
	"fmt"

	filemanager "github.com/nsym-m/simpledb/internal/fileManager"
)

func main() {
	fmt.Println("hello world")
	fm, err := filemanager.NewFileManager("filetest", 400)
	if err != nil {
		panic(err)
	}
	block := filemanager.NewBlockID("testfile", 2)
	page1 := filemanager.NewPage(fm.BlockSize())
	pos1 := 88
	text := "abcdefghijklm"
	page1.SetString(pos1, text)
	size := filemanager.MaxLength(len(text))
	pos2 := pos1 + size
	page1.SetInt(pos2, 345)
	fm.Write(*block, page1)

	page2 := filemanager.NewPage(fm.BlockSize())
	fm.Read(*block, page2)
	fmt.Printf("pos2 offset: %v, contains: %v\n", pos2, page2.GetInt(pos2))
	fmt.Printf("pos1 offset: %v, contains: %v\n", pos1, page2.GetString(pos1))
}
