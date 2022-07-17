package mmap

import (
	"os"

)

// Mmap 利用系统调用中的mmap映射进行文件的读取，在文件读取区间要进行内存保护
func Mmap(fd *os.File, writable bool, size int64) ([]byte, error) {
	return mmap(fd, writable, size)
}

// Munmap 解除文件的mmap映射
func Munmap(b []byte) error {
	return munmap(b)
}

// Madvise 用madvise系统调用给出一些加速内存的建议，其实就是改寄存器的地址
func Madvise(b []byte, readAhead bool) error {
	return madvise(b, readAhead)
}

// Msync 会对已经mmap的文件进行同步操作
func Msync(b []byte) error {
	return msync(b)
}