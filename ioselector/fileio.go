package ioselector

import (
	"os"

)

// FileIOSelector 标准文件IO
type FileIOSelector struct {
	fd *os.File
}

func (fio *FileIOSelector)Write(b []byte, offset int64) (int, error) {
	return fio.fd.WriteAt(b, offset)
}

func (fio *FileIOSelector)Read(b []byte, offset int64) (int, error) {
	return fio.fd.ReadAt(b, offset)
}

func (fio *FileIOSelector) Sync() error {
	return fio.fd.Sync()
}

func (fio *FileIOSelector) Close() error {
	return fio.fd.Close()
}

func (fio *FileIOSelector) Delete() error {
	if err := fio.fd.Close(); err != nil {
		return err
	}
	return os.Remove(fio.fd.Name())
}

// NewFileIOSelector 构造文件IO选择器
func NewFileIOSelector(fileName string, fileSize int64)(IOSelector, error) {
	if fileSize <= 0 {
		return nil, ErrInvalidFileSize
	}
	file, err := openFile(fileName, fileSize)
	if err != nil {
		return nil, err
	}
	return &FileIOSelector{fd: file}, nil
}