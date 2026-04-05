package log

import (
	"iter"
	"sync"

	apperrors "github.com/nsym-m/simpledb/internal/errors"
	"github.com/nsym-m/simpledb/internal/file"
)

type Appender interface {
	Flush(lsn int) error
	Append(logrec []byte) (int, error)
	All() iter.Seq2[[]byte, error]
}

type appender struct {
	blockStore   file.BlockStore
	logFile      string
	logPage      file.Page
	currentBlock *file.BlockID
	latestLSN    int
	lastSavedLSN int
	mu           sync.Mutex
}

func NewAppender(blockStore file.BlockStore, logFile string) (*appender, error) {
	b := make([]byte, blockStore.BlockSize())
	logPage := file.NewPageFromBytes(b)
	logSize, err := blockStore.BlockCount(logFile)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.AppenderInitCode, "appender initialization failed", err)
	}
	lm := &appender{
		blockStore: blockStore,
		logFile:    logFile,
		logPage:    logPage,
	}
	var currentBlock *file.BlockID
	if logSize == 0 {
		currentBlock, err = lm.appendNewBlock()
		if err != nil {
			return nil, apperrors.Wrap(apperrors.AppenderInitCode, "appender initialization failed", err)
		}
	} else {
		currentBlock = file.NewBlockID(logFile, logSize-1)
		if err := blockStore.Read(*currentBlock, logPage); err != nil {
			return nil, apperrors.Wrap(apperrors.AppenderInitCode, "appender initialization failed", err)
		}
	}
	lm.currentBlock = currentBlock

	return lm, nil
}

func (lm *appender) Flush(lsn int) error {
	if lsn < lm.lastSavedLSN {
		return nil
	}
	if err := lm.flush(); err != nil {
		return err
	}
	return nil
}

func (lm *appender) Append(logrec []byte) (int, error) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	boundary := lm.logPage.GetInt(0)
	recsize := len(logrec)
	bytesNeeded := recsize + file.Int32Bytes
	if boundary-bytesNeeded < file.Int32Bytes {
		if err := lm.flush(); err != nil {
			return 0, err
		}
		newBlock, err := lm.appendNewBlock()
		if err != nil {
			return 0, err
		}
		lm.currentBlock = newBlock
		boundary = lm.logPage.GetInt(0)
	}
	recpos := boundary - bytesNeeded
	if err := lm.logPage.SetBytes(recpos, logrec); err != nil {
		return 0, err
	}
	if err := lm.logPage.SetInt(0, recpos); err != nil {
		return 0, err
	}
	lm.latestLSN += 1
	return lm.latestLSN, nil
}

func (lm *appender) All() iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		lm.mu.Lock()
		if err := lm.flush(); err != nil {
			yield(nil, err)
			return
		}
		currentBlock := lm.currentBlock
		lm.mu.Unlock()

		for {
			page := file.NewPage(lm.blockStore.BlockSize())
			if err := lm.blockStore.Read(*currentBlock, page); err != nil {
				yield(nil, err)
				return
			}

			boundary := page.GetInt(0)
			for boundary < lm.blockStore.BlockSize() {
				rec := page.GetBytes(boundary)
				if !yield(rec, nil) {
					return // breakされた場合
				}
				boundary += file.Int32Bytes + len(rec)
			}
			if currentBlock.Number() == 0 {
				return
			}
			currentBlock = file.NewBlockID(lm.logFile, currentBlock.Number()-1)
		}
	}
}

func (lm *appender) appendNewBlock() (*file.BlockID, error) {
	newBlock, err := lm.blockStore.Append(lm.logFile)
	if err != nil {
		return nil, err
	}
	if err := lm.logPage.SetInt(0, lm.blockStore.BlockSize()); err != nil {
		return nil, err
	}
	if err := lm.blockStore.Write(*newBlock, lm.logPage); err != nil {
		return nil, err
	}
	return newBlock, nil
}

func (lm *appender) flush() error {
	if err := lm.blockStore.Write(*lm.currentBlock, lm.logPage); err != nil {
		return err
	}
	lm.lastSavedLSN = lm.latestLSN
	return nil
}
