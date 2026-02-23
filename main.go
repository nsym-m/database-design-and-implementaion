package main

import (
	"fmt"

	"github.com/nsym-m/simpledb/internal/file"
)

func main() {
	fmt.Println("hello world")
	fm, err := file.NewFileManager("filetest", 400)
	if err != nil {
		panic(err)
	}
	block := file.NewBlockID("testfile", 2)
	page1 := file.NewPage(fm.BlockSize())
	pos1 := 88
	text := "abcdefghijklm"
	page1.SetString(pos1, text)
	size := file.MaxLength(len(text))
	pos2 := pos1 + size
	page1.SetInt(pos2, 345)
	fm.Write(*block, page1)

	page2 := file.NewPage(fm.BlockSize())
	fm.Read(*block, page2)
	fmt.Printf("pos2 offset: %v, contains: %v\n", pos2, page2.GetInt(pos2))
	fmt.Printf("pos1 offset: %v, contains: %v\n", pos1, page2.GetString(pos1))
}
