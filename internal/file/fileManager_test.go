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
	if err := page1.SetString(pos1, text); err != nil {
		t.Errorf("SetString(pos1, text) pos1: %v, text: %v, error: %v", pos1, text, err)
	}
	size := file.MaxLength(len(text))
	pos2 := pos1 + size

	if err := page1.SetInt(pos2, 345); err != nil {
		t.Errorf("SetInt(pos2, 345) pos2: %v, error: %v", pos2, err)
	}
	if err := fm.Write(*block, *page1); err != nil {
		t.Errorf("Write(*block, *page1) block: %v, page1: %v, error: %v", *block, *page1, err)
	}
	page2 := file.NewPage(fm.BlockSize())
	if err := fm.Read(*block, *page2); err != nil {
		t.Errorf("Read(*block, *page2) block: %v, page2: %v, error: %v", *block, *page2, err)
	}

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
