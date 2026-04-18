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

type Pool interface {
	UnPin(buff Buffer)
	Pin(block *file.BlockID) (Buffer, error)
	Available() int
	FlushAll(txnum int)
}

type pool struct {
	bufferPool   []*buffer
	numAvailable int

	cond sync.Cond
	ctx  context.Context
}

func NewPool(blockStore file.BlockStore, appender log.Appender, numBuffs int) Pool {
	buffers := make([]*buffer, numBuffs)
	for i := range numBuffs {
		buffers[i] = &buffer{
			blockStore: blockStore,
			appender:   appender,
			contents:   file.NewPage(blockStore.BlockSize()),
			txnum:      -1,
		}
	}

	return &pool{
		bufferPool:   buffers,
		numAvailable: numBuffs,
		cond: sync.Cond{
			L: &sync.Mutex{},
		},
	}
}

func (p *pool) UnPin(buff Buffer) {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()

	buff.UnPin()
	if !buff.IsPinned() {
		p.numAvailable++
		p.cond.Broadcast()
	}
}

func (p *pool) Pin(block *file.BlockID) (Buffer, error) {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()

	now := time.Now()
	start := now.UnixMilli()
	dead := now.Add(time.Duration(MaxTime) * time.Millisecond)
	buff := p.tryToPin(block)
	for buff == nil && !p.waitingTooLong(start) {
		syncutil.WaitWithDeadline(&p.cond, dead)
		buff = p.tryToPin(block)
	}
	if buff == nil {
		return buff, apperrors.New(apperrors.BufferAbortCode, "buffer abort")
	}
	return buff, nil
}

func (p *pool) Available() int {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()
	return p.numAvailable
}

func (p *pool) FlushAll(txnum int) {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()
	for _, buff := range p.bufferPool {
		if buff.ModifyingTx() == txnum {
			p.numAvailable++
			p.cond.Broadcast()
		}
	}
}

func (p *pool) waitingTooLong(start int64) bool {
	return time.Now().UnixMilli()-start > MaxTime
}

func (p *pool) tryToPin(block *file.BlockID) *buffer {
	buff := p.findExistingBuffer(block)
	if buff == nil {
		buff = p.chooseUnPinnedBuffer()
		if buff == nil {
			return nil
		}
		buff.AssignToBlock(block)
	}
	if !buff.IsPinned() {
		p.numAvailable--
	}
	buff.Pin()
	return buff
}

func (p *pool) findExistingBuffer(block *file.BlockID) *buffer {
	for _, buf := range p.bufferPool {
		b := buf.Block()
		if b != nil && b == block {
			return buf
		}
	}
	return nil
}

func (p *pool) chooseUnPinnedBuffer() *buffer {
	for _, buf := range p.bufferPool {
		if !buf.IsPinned() {
			return buf
		}
	}
	return nil
}
