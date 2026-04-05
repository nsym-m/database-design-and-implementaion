package buffer

import (
	"sync"
	"time"

	apperrors "github.com/nsym-m/simpledb/internal/errors"
	"github.com/nsym-m/simpledb/internal/file"
	"github.com/nsym-m/simpledb/internal/log"
)

const MaxTime int64 = 1000

type BufferMgr interface {
}

type bufferMgr struct {
	bufferPool   []*buffer
	numAvailable int
	mu           sync.Mutex
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
	return bufferMgr{
		bufferPool:   pool,
		numAvailable: numBuffs,
	}
}

func (bm *bufferMgr) UnPin(buff Buffer) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	buff.UnPin()
	if !buff.IsPinned() {
		bm.numAvailable++
		// bm.NotifyAll()
	}
}

func (bm *bufferMgr) Pin(block *file.BlockID) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	start := time.Now().UnixMilli()
	buff := bm.tryToPin(block)
	for buff == nil && !bm.waitingTooLong(start) {
		// wait()
		buff = bm.tryToPin(block)
	}
	if buff == nil {
		return apperrors.ErrBufferAbort
	}
	return nil
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
