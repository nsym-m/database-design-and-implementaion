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
	buffer.NewBuffer(bs, ap)
}
