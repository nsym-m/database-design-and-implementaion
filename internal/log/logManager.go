package log

import (
	"fmt"
	"iter"
	"sync"

	"github.com/nsym-m/simpledb/internal/file"
)

const Bytes = 4

type LogManager struct {
	fileManager  *file.FileManager
	logFile      string
	logPage      *file.Page
	currentBlock *file.BlockID
	latestLSN    int
	lastSavedLSN int
	mu           sync.Mutex
}

func NewLogManager(fileManager *file.FileManager, logFile string) (*LogManager, error) {
	b := make([]byte, fileManager.BlockSize())
	logPage := file.NewPageFromBytes(b)
	logSize, err := fileManager.BlockCount(logFile)
	if err != nil {
		return nil, fmt.Errorf("NewLogManager error: %w", err)
	}
	lm := &LogManager{
		fileManager: fileManager,
		logFile:     logFile,
		logPage:     logPage,
	}
	var currentBlock *file.BlockID
	if logSize == 0 {
		currentBlock, err = lm.appendNewBlock()
		if err != nil {
			return nil, fmt.Errorf("NewLogManager error: %w", err)
		}
	} else {
		currentBlock = file.NewBlockID(logFile, logSize-1)
		if err := fileManager.Read(*currentBlock, *logPage); err != nil {
			return nil, fmt.Errorf("NewLogManager error: %w", err)
		}
	}
	lm.currentBlock = currentBlock

	return lm, nil
}

func (lm *LogManager) Flush(lsn int) error {
	if lsn < lm.lastSavedLSN {
		return nil
	}
	if err := lm.flush(); err != nil {
		return err
	}
	return nil
}

func (lm *LogManager) Append(logrec []byte) (int, error) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	boundary := lm.logPage.GetInt(0)
	recsize := len(logrec)
	bytesNeeded := recsize + Bytes
	if boundary-bytesNeeded < Bytes {
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

func (lm *LogManager) All() iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		lm.mu.Lock()
		if err := lm.flush(); err != nil {
			yield(nil, err)
			return
		}
		currentBlock := lm.currentBlock
		lm.mu.Unlock()

		for {
			page := file.NewPage(lm.fileManager.BlockSize())
			if err := lm.fileManager.Read(*currentBlock, *page); err != nil {
				yield(nil, err)
				return
			}

			boundary := page.GetInt(0)
			for boundary < lm.fileManager.BlockSize() {
				rec := page.GetBytes(boundary)
				if !yield(rec, nil) {
					return // breakされた場合
				}
				boundary += Bytes + len(rec)
			}
			if currentBlock.Number() == 0 {
				return
			}
			currentBlock = file.NewBlockID(lm.logFile, currentBlock.Number()-1)
		}
	}
}

func (lm *LogManager) appendNewBlock() (*file.BlockID, error) {
	newBlock, err := lm.fileManager.Append(lm.logFile)
	if err != nil {
		return nil, err
	}
	if err := lm.logPage.SetInt(0, lm.fileManager.BlockSize()); err != nil {
		return nil, err
	}
	if err := lm.fileManager.Write(*newBlock, *lm.logPage); err != nil {
		return nil, err
	}
	return newBlock, nil
}

func (lm *LogManager) flush() error {
	if err := lm.fileManager.Write(*lm.currentBlock, *lm.logPage); err != nil {
		return err
	}
	lm.lastSavedLSN = lm.latestLSN
	return nil
}
