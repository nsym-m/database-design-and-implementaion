package buffer

import (
	"context"
	"sync"
	"time"

	apperrors "github.com/nsym-m/simpledb/internal/errors"
	"github.com/nsym-m/simpledb/internal/file"
	"github.com/nsym-m/simpledb/internal/log"
	"github.com/nsym-m/simpledb/internal/syncutil"
)

const MaxTime int64 = 1000

type BufferMgr interface {
	UnPin(buff Buffer)
	Pin(block *file.BlockID) (Buffer, error)
	Available() int
	FlushAll(txnum int)
}

type bufferMgr struct {
	bufferPool   []*buffer
	numAvailable int

	cond sync.Cond
	ctx  context.Context
}

func NewBufferMgr(blockStore file.BlockStore, appender log.Appender, numBuffs int) BufferMgr {
	pool := make([]*buffer, numBuffs)
	for i := range numBuffs {
		pool[i] = &buffer{
			blockStore: blockStore,
			appender:   appender,
			contents:   file.NewPage(blockStore.BlockSize()),
		}
	}

	return &bufferMgr{
		bufferPool:   pool,
		numAvailable: numBuffs,
		cond: sync.Cond{
			L: &sync.Mutex{},
		},
	}
}

func (bm *bufferMgr) UnPin(buff Buffer) {
	bm.cond.L.Lock()
	defer bm.cond.L.Unlock()

	buff.UnPin()
	if !buff.IsPinned() {
		bm.numAvailable++
		bm.cond.Broadcast()
	}
}

func (bm *bufferMgr) Pin(block *file.BlockID) (Buffer, error) {
	bm.cond.L.Lock()
	defer bm.cond.L.Unlock()

	now := time.Now()
	start := now.UnixMilli()
	buff := bm.tryToPin(block)
	for buff == nil && !bm.waitingTooLong(start) {
		syncutil.WaitWithDeadline(&bm.cond, now)
		buff = bm.tryToPin(block)
	}
	if buff == nil {
		return buff, apperrors.New(apperrors.BufferAbortCode, "buffer abort")
	}
	return buff, nil
}

func (bm *bufferMgr) Available() int {
	bm.cond.L.Lock()
	defer bm.cond.L.Unlock()
	return bm.numAvailable
}

func (bm *bufferMgr) FlushAll(txnum int) {
	bm.cond.L.Lock()
	defer bm.cond.L.Unlock()
	for _, buff := range bm.bufferPool {
		if buff.ModifyingTx() == txnum {
			bm.numAvailable++
			bm.cond.Broadcast()
		}
	}
}

func (bm *bufferMgr) waitingTooLong(start int64) bool {
	return time.Now().UnixMilli()-start > MaxTime
}

func (bm *bufferMgr) tryToPin(block *file.BlockID) *buffer {
	buff := bm.findExistingBuffer(block)
	if buff == nil {
		buff = bm.chooseUnPinnedBuffer()
		if buff != nil {
			buff.AssignToBlock(block)
		}
	}
	if !buff.IsPinned() {
		bm.numAvailable--
	}
	buff.Pin()
	return buff
}

func (bm *bufferMgr) findExistingBuffer(block *file.BlockID) *buffer {
	for _, buffer := range bm.bufferPool {
		b := buffer.Block()
		if b != nil && b == block {
			return buffer
		}
	}
	return nil
}

func (bm *bufferMgr) chooseUnPinnedBuffer() *buffer {
	for _, buffer := range bm.bufferPool {
		if !buffer.IsPinned() {
			return buffer
		}
	}
	return nil
}
