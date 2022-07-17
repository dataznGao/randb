package ioselector

import (
	"io"
	"os"
	"randb/mmap"
)

type MMapSelector struct {
	fd   	*os.File
	// buf mmap的缓冲
	buf  	[]byte
	bufLen 	int64
}

// Write 往selector里面写数据，其实就是把[]byte进行操作
func (lm *MMapSelector) Write(b []byte, offset int64)(int, error) {
	length := int64(len(b))
	if length <= 0 {
		return 0, nil
	}
	if offset < 0 || offset+length > lm.bufLen {
		return 0, io.EOF
	}
	return copy(lm.buf[offset:], b), nil
}

func (lm *MMapSelector) Read(b []byte, offset int64)(int, error) {
	if offset < 0 || offset + int64(len(b)) >= lm.bufLen || offset >= lm.bufLen {
		return 0, io.EOF
	}
	return copy(b, lm.buf[offset:]), nil
}

func (lm *MMapSelector) Sync() error {
	return mmap.Msync(lm.buf)
}

func (lm *MMapSelector) Close() error {
	// 1. 同步落盘
	if err := mmap.Msync(lm.buf); err != nil {
		return err
	}
	// 2. 解除映射
	if err := mmap.Munmap(lm.buf); err != nil {
		return err
	}
	return nil
}

func (lm *MMapSelector) Delete() error {
	// 1. 解除映射
	if err := mmap.Munmap(lm.buf); err != nil {
		return err
	}
	// 2. selector buf置空
	lm.buf = nil
	// 3. 长度截断到0
	if err := lm.fd.Truncate(0); err != nil {
		return err
	}
	// 4. 关闭文件描述符
	if err := lm.fd.Close(); err != nil {
		return err
	}
	// 5. 删除文件
	return os.Remove(lm.fd.Name())
}


// NewMMapSelector 构造新的mmap selector
func NewMMapSelector(fileName string, fileSize int64) (IOSelector, error) {
	if fileSize < 0 {
		return nil, ErrInvalidFileSize
	}
	// 打开文件
	file, err := openFile(fileName, fileSize)
	// 进行mmap映射，映射到一个buffer中
	if err != nil {
		return nil, err
	}

	buf, err := mmap.Mmap(file, true, fileSize)
	if err != nil {
		return nil, err
	}
	return &MMapSelector{fd: file, buf: buf, bufLen: int64(len(buf))}, nil
}