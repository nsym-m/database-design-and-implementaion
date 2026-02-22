package filemanager

// BlockID 複数ブロックに分割されたDBのファイルをブロックごとに識別する構造体
type BlockID struct {
	fileName string
	blockNum int
}

func NewBlockID(fileName string, blockNum int) *BlockID {
	return &BlockID{
		fileName: fileName,
		blockNum: blockNum,
	}
}

func (b BlockID) FileName() string {
	return b.fileName
}

func (b BlockID) Number() int {
	return b.blockNum
}
