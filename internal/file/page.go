package file

import (
	"encoding/binary"
	"fmt"
	"math"
	"unicode/utf8"
)

// Page ブロックサイズのメモリ領域
type Page struct {
	byteBuffer []byte
}

func NewPage(blockSize int) *Page {
	return &Page{
		byteBuffer: make([]byte, blockSize),
	}
}

func NewPageFromBytes(bytes []byte) *Page {
	fmt.Printf("Page: %+v\n", &Page{
		byteBuffer: bytes,
	})

	return &Page{
		byteBuffer: bytes,
	}
}

func (p Page) SetInt(offset, i int) error {
	if i > math.MaxInt32 || i < math.MinInt32 {
		return fmt.Errorf("SetInt: value %d overflows int32", i)
	}
	//nolint:gosec // int32の範囲チェックをしているので問題なし
	binary.BigEndian.PutUint32(p.byteBuffer[offset:], uint32(int32(i)))
	return nil
}

func (p Page) GetInt(offset int) int {
	//nolint:gosec // uint32->int32の変換は範囲外にならないので問題なし
	return int(int32(binary.BigEndian.Uint32(p.byteBuffer[offset:])))
}

func (p Page) SetBytes(offset int, b []byte) error {
	length := len(b)
	if length > math.MaxInt32 || length < math.MinInt32 {
		return fmt.Errorf("SetBytes: len(b) %d overflows int32", length)
	}
	// 長さを先に書く
	binary.BigEndian.PutUint32(p.byteBuffer[offset:], uint32(int32(length)))
	// バイト列を書く
	copy(p.byteBuffer[offset+4:], b)
	return nil
}

func (p Page) GetBytes(offset int) []byte {
	// 長さを読む
	//nolint:gosec // uint32->int32の変換は範囲外にならないので問題なし
	length := int(int32(binary.BigEndian.Uint32(p.byteBuffer[offset:])))
	// 長さ分の領域を確保してコピー
	b := make([]byte, length)
	copy(b, p.byteBuffer[offset+4:])
	return b
}

func (p Page) SetString(offset int, s string) error {
	if err := p.SetBytes(offset, []byte(s)); err != nil {
		return err
	}
	return nil
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
