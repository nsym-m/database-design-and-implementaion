package buffer_test

import (
	"testing"

	"github.com/nsym-m/simpledb/internal/buffer"
	"github.com/nsym-m/simpledb/internal/file"
	"github.com/nsym-m/simpledb/internal/log"
)

func TestBuffer(t *testing.T) {

	tmpDir := t.TempDir()
	bs, err := file.NewBlockStore(tmpDir, 400)
	if err != nil {
		t.Fatal(err)
	}

	ap, err := log.NewAppender(bs, "test")
	if err != nil {
		t.Fatal(err)
	}

	pool := buffer.NewPool(bs, ap, 3)

	buff1, err := pool.Pin(file.NewBlockID("testfile", 1))
	if err != nil {
		t.Fatal(err)
	}
	p := buff1.Contents()
	n := p.GetInt(80)
	t.Logf("n1 %d\n", n)
	p.SetInt(80, n+1)
	t.Logf("n2 %d\n", p.GetInt(80))
	buff1.SetModified(1, 0)
	t.Logf("new value %d\n", n+1)
	pool.UnPin(buff1)
	buff2, err := pool.Pin(file.NewBlockID("testfile", 2))
	if err != nil {
		t.Fatal(err)
	}
	// buff3, err := bm.Pin(file.NewBlockID("testfile", 3))
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// buff4, err := bm.Pin(file.NewBlockID("testfile", 4))
	// if err != nil {
	// 	t.Fatal(err)
	// }

	pool.UnPin(buff2)
	buff2, err = pool.Pin(file.NewBlockID("testfile", 1))
	if err != nil {
		t.Fatal(err)
	}
	p2 := buff2.Contents()
	p2.SetInt(80, 9999)
	buff2.SetModified(1, 0)
	pool.UnPin(buff2)
}
