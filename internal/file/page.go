package file

import (
	"encoding/binary"
	"unicode/utf8"
)

// Page ブロックサイズのメモリ領域
type Page struct {
	byteBuffer []byte
}

func NewPage(blockSize int) Page {
	return Page{
		byteBuffer: make([]byte, blockSize),
	}
}

func NewPageFromBytes(bytes []byte) Page {
	return Page{
		byteBuffer: bytes,
	}
}

func (p Page) SetInt(offset, i int) {
	binary.BigEndian.PutUint32(p.byteBuffer[offset:], uint32(int32(i)))
}

func (p Page) GetInt(offset int) int {
	return int(int32(binary.BigEndian.Uint32(p.byteBuffer[offset:])))
}

func (p Page) SetBytes(offset int, b []byte) {
	// 長さを先に書く
	binary.BigEndian.PutUint32(p.byteBuffer[offset:], uint32(int32(len(b))))
	// バイト列を書く
	copy(p.byteBuffer[offset+4:], b)
}

func (p Page) GetBytes(offset int) []byte {
	// 長さを読む
	length := int(int32(binary.BigEndian.Uint32(p.byteBuffer[offset:])))
	// 長さ分の領域を確保してコピー
	b := make([]byte, length)
	copy(b, p.byteBuffer[offset+4:])
	return b
}

func (p Page) SetString(offset int, s string) {
	p.SetBytes(offset, []byte(s))
}

func (p Page) GetString(offset int) string {
	return string(p.GetBytes(offset))
}

func MaxLength(strlen int) int {
	return 4 + (strlen * utf8.UTFMax)
}

func (p Page) contents() []byte {
	return p.byteBuffer
}
