package file

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	dirPerm fs.FileMode = 0700
)

const (
	Int32Bytes int = 4
)

// BlockStore DBのデータを管理するファイルを操作する
type BlockStore interface {
	Read(block BlockID, page Page) error
	Write(block BlockID, page Page) error
	Append(fileName string) (*BlockID, error)
	BlockCount(fileName string) (int, error)
	BlockSize() int
}

type blockStore struct {
	dbPath    string
	blockSize int
	isNew     bool
	openFiles map[string]*os.File
	mu        sync.Mutex
}

func NewBlockStore(dbPath string, blockSize int) (*blockStore, error) {
	isNew := false
	// pathがなければ作成
	if _, err := os.Stat(dbPath); err != nil {
		if err := os.MkdirAll(dbPath, dirPerm); err != nil {
			return nil, fmt.Errorf("MkdirAll error: %w", err)
		}
		isNew = true
	}
	dirs, err := os.ReadDir(dbPath)
	if err != nil {
		return nil, fmt.Errorf("NewBlockStore error: %w", err)
	}
	// 一時DBファイルが存在していたら削除
	for _, file := range dirs {
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(file.Name(), "temp") {
			os.Remove(filepath.Join(dbPath, file.Name()))
		}
	}
	return &blockStore{
		dbPath:    dbPath,
		blockSize: blockSize,
		isNew:     isNew,
		openFiles: make(map[string]*os.File),
	}, nil
}

func (bs *blockStore) Read(block BlockID, page Page) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	file, err := bs.file(block.FileName())
	if err != nil {
		return fmt.Errorf("Read error: %w", err)
	}
	offset := int64(block.Number()) * int64(bs.blockSize)
	if _, err := file.ReadAt(page.Contents(), offset); err != nil {
		return fmt.Errorf("Read error: %w", err)
	}
	return nil
}

func (bs *blockStore) Write(block BlockID, page Page) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	file, err := bs.file(block.FileName())
	if err != nil {
		return fmt.Errorf("Write error: %w", err)
	}
	offset := int64(block.Number()) * int64(bs.blockSize)
	if _, err := file.WriteAt(page.Contents(), offset); err != nil {
		return fmt.Errorf("Write error: %w", err)
	}
	return nil
}

func (bs *blockStore) Append(fileName string) (*BlockID, error) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	newBlockNum, err := bs.BlockCount(fileName)
	if err != nil {
		return nil, fmt.Errorf("BlockCount error: %w", err)
	}
	blockID := NewBlockID(fileName, newBlockNum)
	b := make([]byte, newBlockNum)
	file, err := bs.file(fileName)
	if err != nil {
		return nil, fmt.Errorf("Append error: %w", err)
	}
	offset := int64(blockID.Number()) * int64(bs.blockSize)
	if _, err := file.WriteAt(b, offset); err != nil {
		return nil, fmt.Errorf("Append error: %w", err)
	}
	return blockID, nil
}

func (bs *blockStore) BlockSize() int {
	return bs.blockSize
}

func (bs *blockStore) IsNew() bool {
	return bs.isNew
}

func (bs *blockStore) BlockCount(fileName string) (int, error) {
	file, err := bs.file(fileName)
	if err != nil {
		return 0, fmt.Errorf("BlockCount error: %w", err)
	}
	info, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("BlockCount error: %w", err)
	}
	return int(info.Size()) / bs.blockSize, nil
}

func (bs *blockStore) file(fileName string) (*os.File, error) {
	f, ok := bs.openFiles[fileName]
	if ok {
		return f, nil
	}
	f, err := os.OpenFile(filepath.Join(bs.dbPath, fileName), os.O_RDWR|os.O_CREATE, dirPerm)
	if err != nil {
		return nil, fmt.Errorf("getFile error: %w", err)
	}
	bs.openFiles[fileName] = f
	return f, nil
}
