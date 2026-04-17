package buffer

import (
	"github.com/nsym-m/simpledb/internal/file"
	"github.com/nsym-m/simpledb/internal/log"
)

type Buffer interface {
	Contents() file.Page
	Block() *file.BlockID
	SetModified(txmun int, lsn int)
	IsPinned() bool
	ModifyingTx() int
	AssignToBlock(block *file.BlockID)
	Flush()
	Pin()
	UnPin()
}

type buffer struct {
	blockStore file.BlockStore
	appender   log.Appender
	contents   file.Page
	block      *file.BlockID
	pins       int
	txnum      int
	lsn        int
}

func (b *buffer) Contents() file.Page {
	return b.contents
}

func (b *buffer) Block() *file.BlockID {
	return b.block
}

func (b *buffer) SetModified(txmun int, lsn int) {
	b.txnum = txmun
	if lsn >= 0 {
		b.lsn = lsn
	}
}

func (b *buffer) IsPinned() bool {
	return b.pins > 0
}

func (b *buffer) ModifyingTx() int {
	return b.txnum
}

func (b *buffer) AssignToBlock(block *file.BlockID) {
	b.Flush()
	b.block = block
	b.blockStore.Read(*b.block, b.contents)
	b.pins = 0
}

func (b *buffer) Flush() {
	if b.txnum >= 0 {
		b.appender.Flush(b.lsn)
		// TODO: *b.block で panic
		b.blockStore.Write(*b.block, b.contents)
		b.txnum = -1
	}
}

func (b *buffer) Pin() {
	b.pins++
}

func (b *buffer) UnPin() {
	b.pins--
}
