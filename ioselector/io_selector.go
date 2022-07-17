package ioselector

import (
	"errors"
	"os"
)

// ErrInvalidFileSize 错误定义：无效的文件大小
var ErrInvalidFileSize = errors.New("file size can`t be zero or negative")

// FilePerm 新建的文件的授权
const FilePerm = 0644

// IOSelector 读写选择器是用来进行文件读写或者mmap映射的
type IOSelector interface {
	// Write 将byte数组追加到offset的位置
	// 它返回的是写入的字节数
	Write(b []byte, offset int64) (int, error)

	// Read 从offset的位置开始读，读的数据放在b中
	Read(b []byte, offset int64) (int, error)

	// Sync 提交文件的当前内容到磁盘中
	Sync() error

	// Close 关闭文件
	// 当文件已经关闭时，返回error
	Close() error

	// Delete 删除文件
	// 删除文件前必须关闭文件
	Delete() error
}

// 获取文件描述符，如果文件比这个size小，就截断，也就是扩大到这个大小
func openFile(fileName string, fileSize int64) (*os.File, error) {
	fd, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, FilePerm)
	if err != nil {
		return fd, err
	}

	stat, err := fd.Stat()
	if err != nil {
		return nil, err
	}

	if stat.Size() < fileSize {
		if err := fd.Truncate(fileSize); err != nil {
			return nil, err
		}
	}
	return fd, nil
}
