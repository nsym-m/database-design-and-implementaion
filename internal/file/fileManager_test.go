package file_test

import (
	"testing"

	"github.com/nsym-m/simpledb/internal/file"
)

func TestFileManager(t *testing.T) {

	tmpDir := t.TempDir()
	fm, err := file.NewFileManager(tmpDir, 400)
	if err != nil {
		t.Fatal(err)
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

	// test
	got := page2.GetInt(pos2)
	if got != 345 {
		t.Errorf("GetInt(%d) = %d, want 345", pos2, got)
	}
	gotStr := page2.GetString(pos1)
	if gotStr != text {
		t.Errorf("GetString(%d) = %q, want %q", pos1, gotStr, text)
	}
}
