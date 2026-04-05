package file

import (
	"encoding/binary"
	"fmt"
	"math"
	"unicode/utf8"

	apperrors "github.com/nsym-m/simpledb/internal/errors"
)

// Page ブロックサイズのメモリ領域
type Page interface {
	SetInt(offset, i int) error
	GetInt(offset int) int
	SetBytes(offset int, b []byte) error
	GetBytes(offset int) []byte
	SetString(offset int, s string) error
	GetString(offset int) string
	Contents() []byte
}

type page struct {
	byteBuffer []byte
}

func NewPage(blockSize int) *page {
	return &page{
		byteBuffer: make([]byte, blockSize),
	}
}

func NewPageFromBytes(bytes []byte) *page {
	return &page{
		byteBuffer: bytes,
	}
}

func (p *page) SetInt(offset, i int) error {
	if i > math.MaxInt32 || i < math.MinInt32 {
		return apperrors.New(apperrors.IntOverflowCode, fmt.Sprintf("SetInt: value %d overflows int32", i))
	}
	//nolint:gosec // int32の範囲チェックをしているので問題なし
	binary.BigEndian.PutUint32(p.byteBuffer[offset:], uint32(int32(i)))
	return nil
}

func (p *page) GetInt(offset int) int {
	//nolint:gosec // uint32->int32の変換は範囲外にならないので問題なし
	return int(int32(binary.BigEndian.Uint32(p.byteBuffer[offset:])))
}

func (p *page) SetBytes(offset int, b []byte) error {
	length := len(b)
	if length > math.MaxInt32 || length < math.MinInt32 {
		return apperrors.New(apperrors.BytesOverflowCode, fmt.Sprintf("SetBytes: len(b) %d overflows int32", length))
	}
	// 長さを先に書く
	binary.BigEndian.PutUint32(p.byteBuffer[offset:], uint32(int32(length)))
	// バイト列を書く
	copy(p.byteBuffer[offset+Int32Bytes:], b)
	return nil
}

func (p *page) GetBytes(offset int) []byte {
	// 長さを読む
	//nolint:gosec // uint32->int32の変換は範囲外にならないので問題なし
	length := int(int32(binary.BigEndian.Uint32(p.byteBuffer[offset:])))
	// 長さ分の領域を確保してコピー
	b := make([]byte, length)
	copy(b, p.byteBuffer[offset+Int32Bytes:])
	return b
}

func (p *page) SetString(offset int, s string) error {
	if err := p.SetBytes(offset, []byte(s)); err != nil {
		return err
	}
	return nil
}

func (p *page) GetString(offset int) string {
	return string(p.GetBytes(offset))
}

func MaxLength(strlen int) int {
	return Int32Bytes + (strlen * utf8.UTFMax)
}

func (p *page) Contents() []byte {
	return p.byteBuffer
}
