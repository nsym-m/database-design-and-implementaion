package errors

import (
	"errors"
	"fmt"
)

// AppError は構造化されたエラー型です。
// errors.Is() / errors.As() / errors.AsType() に対応しています。
type AppError struct {
	Code ErrorCode
	Msg  string
	wrap error
}

func (e *AppError) Error() string {
	if e.wrap != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Msg, e.wrap)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Msg)
}

// Unwrap はエラーチェーンの探索を可能にします。
func (e *AppError) Unwrap() error {
	return e.wrap
}

// New は新しい AppError を生成します。
func New(code ErrorCode, msg string) *AppError {
	return &AppError{Code: code, Msg: msg}
}

// Wrap は既存のエラーをラップした AppError を生成します。
func Wrap(code ErrorCode, msg string, err error) *AppError {
	return &AppError{Code: code, Msg: msg, wrap: err}
}

// IsCode はエラーチェーンに指定した ErrorCode を持つ AppError が含まれるか検査します。
func IsCode(err error, code ErrorCode) bool {
	e, ok := errors.AsType[*AppError](err)
	return ok && e.Code == code
}
