package utils

import "fmt"

// ErrType 定义错误类型
type ErrType int

const (
	ErrFileAccess ErrType = iota
	ErrMetadata
	ErrCopy
	ErrInvalidFormat
)

// MediaError 自定义错误类型
type MediaError struct {
	Type    ErrType
	Path    string
	Message string
	Err     error
}

func (e *MediaError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Message, e.Path, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Message, e.Path)
}

// NewFileAccessError 创建文件访问错误
func NewFileAccessError(path string, err error) error {
	return &MediaError{
		Type:    ErrFileAccess,
		Path:    path,
		Message: "无法访问文件",
		Err:     err,
	}
}

// NewMetadataError 创建元数据读取错误
func NewMetadataError(path string, err error) error {
	return &MediaError{
		Type:    ErrMetadata,
		Path:    path,
		Message: "无法读取元数据",
		Err:     err,
	}
}

// NewCopyError 创建文件复制错误
func NewCopyError(path string, err error) error {
	return &MediaError{
		Type:    ErrCopy,
		Path:    path,
		Message: "复制文件失败",
		Err:     err,
	}
}

// NewInvalidFormatError 创建无效格式错误
func NewInvalidFormatError(path string) error {
	return &MediaError{
		Type:    ErrInvalidFormat,
		Path:    path,
		Message: "不支持的文件格式",
	}
}
