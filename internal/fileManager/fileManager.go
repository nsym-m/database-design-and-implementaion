package filemanager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// FileManager DBのデータを管理するファイルを操作する
type FileManager struct {
	dbPath    string
	blockSize int
	isNew     bool
	openFiles map[string]*os.File
	mu        sync.Mutex
}

func NewFileManager(dbPath string, blockSize int) (*FileManager, error) {
	isNew := false
	// pathがなければ作成
	if _, err := os.Stat(dbPath); err != nil {
		os.MkdirAll(dbPath, 0755)
		isNew = true
	}
	dirs, err := os.ReadDir(dbPath)
	if err != nil {
		return nil, fmt.Errorf("NewFileManager error: %w", err)
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
	return &FileManager{
		dbPath:    dbPath,
		blockSize: blockSize,
		isNew:     isNew,
		openFiles: make(map[string]*os.File),
	}, nil
}

func (fm *FileManager) Read(block BlockID, page Page) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	file, err := fm.getFile(block.FileName())
	if err != nil {
		return fmt.Errorf("Read error: %w", err)
	}
	offset := int64(block.Number()) * int64(fm.blockSize)
	if _, err := file.ReadAt(page.contents(), offset); err != nil {
		return fmt.Errorf("Read error: %w", err)
	}
	return nil
}

func (fm *FileManager) Write(block BlockID, page Page) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	file, err := fm.getFile(block.FileName())
	if err != nil {
		return fmt.Errorf("Write error: %w", err)
	}
	offset := int64(block.Number()) * int64(fm.blockSize)
	if _, err := file.WriteAt(page.contents(), offset); err != nil {
		return fmt.Errorf("Write error: %w", err)
	}
	return nil
}

func (fm *FileManager) Append(fileName string) (*BlockID, error) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	newBlockNum := len(fileName)
	blockID := NewBlockID(fileName, newBlockNum)
	b := make([]byte, newBlockNum)
	file, err := fm.getFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("Append error: %w", err)
	}
	offset := int64(blockID.Number()) * int64(fm.blockSize)
	if _, err := file.WriteAt(b, offset); err != nil {
		return nil, fmt.Errorf("Append error: %w", err)
	}
	return blockID, nil
}

func (fm *FileManager) BlockSize() int {
	return fm.blockSize
}

func (fm *FileManager) getFile(fileName string) (*os.File, error) {
	f, ok := fm.openFiles[fileName]
	if ok {
		return f, nil
	}
	f, err := os.OpenFile(filepath.Join(fm.dbPath, fileName), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("getFile error: %w", err)
	}
	fm.openFiles[fileName] = f
	return f, nil
}
