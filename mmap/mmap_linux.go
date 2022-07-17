package mmap

import (
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

// mmap 利用系统调用中的mmap映射进行文件的读取，在文件读取区间要进行内存保护
func mmap(fd *os.File, writable bool, size int64) ([]byte, error) {
	// 只读
	mtype := unix.PROT_READ
	if writable {
		mtype |= unix.PROT_WRITE
	}

	return unix.Mmap(int(fd.Fd()), 0, int(size), mtype, unix.MAP_SHARED)
}

// munmap 解除先前的映射
func munmap(data []byte) error {
	if len(data) == 0 || len(data) != cap(data) {
		return unix.EINVAL
	}
	_, _, errno := unix.Syscall(
		unix.SYS_MUNMAP,
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		0,
	)
	if errno != 0 {
		return errno
	}
	return nil
}

// madvise 当做mmap时，madvise可以给出内存使用的建议，如果是有序的可以开启readAhead预读
func madvise(b []byte, readAhead bool) error {
	flags := unix.MADV_NORMAL
	if !readAhead {
		flags = unix.MADV_RANDOM
	}
	return unix.Madvise(b, flags)
}

// msync 同步mmap数组到磁盘
func msync(b []byte) error {
	return unix.Msync(b, unix.MS_SYNC)
}