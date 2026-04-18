package buffer_test

import (
	"testing"

	"github.com/nsym-m/simpledb/internal/buffer"
	"github.com/nsym-m/simpledb/internal/file"
	"github.com/nsym-m/simpledb/internal/log"
)

func TestBufferMgr(t *testing.T) {

	tmpDir := t.TempDir()
	bs, err := file.NewBlockStore(tmpDir, 400)
	if err != nil {
		t.Fatal(err)
	}

	ap, err := log.NewAppender(bs, "test")
	if err != nil {
		t.Fatal(err)
	}

	buff := make([]buffer.Buffer, 6)

	bm := buffer.NewBufferMgr(bs, ap, 3)

	buff[0], err = bm.Pin(file.NewBlockID("testfile", 0))
	if err != nil {
		t.Fatal(err)
	}
	buff[1], err = bm.Pin(file.NewBlockID("testfile", 1))
	if err != nil {
		t.Fatal(err)
	}
	buff[2], err = bm.Pin(file.NewBlockID("testfile", 2))
	if err != nil {
		t.Fatal(err)
	}
	bm.UnPin(buff[1])
	buff[1] = nil
	buff[3], err = bm.Pin(file.NewBlockID("testfile", 3))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("available buffers: %v\n", bm.Available())
	if bm.Available() != 3 {
		t.Errorf("available buffers want 3, but %d", bm.Available())
	}

	t.Log("Attempting to pin block 3...")
	buff[5], err = bm.Pin(file.NewBlockID("testfile", 3))
	if err != nil {
		t.Log("Exception: No available buffers")
	}
	bm.UnPin(buff[2])
	buff[2] = nil
	buff[5], err = bm.Pin(file.NewBlockID("testfile", 3))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Final Buffer Allocation:")
	count := len(buff)
	for i := 0; i < count; i++ {
		b := buff[i]
		if b != nil {
			t.Logf("Buff[%d] pinned to block %+v\n", i, b.Block())
		}
	}
}
