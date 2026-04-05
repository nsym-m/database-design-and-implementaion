package errors

// ErrorCode はエラーの種別を表す識別子です。
type ErrorCode string

const (
	BufferAbortCode   ErrorCode = "E0001"
	IntOverflowCode   ErrorCode = "E0002"
	BytesOverflowCode ErrorCode = "E0003"
	AppenderInitCode  ErrorCode = "E0004"
	BlockStoreIOCode  ErrorCode = "E0005"
)
