//go:build !windows && !plan9
// +build !windows,!plan9

package flock

import (
	"os"
	"syscall"
)

//FileLockGuard 持有了一个在目录上的文件锁
type FileLockGuard struct {
	//在目录上的文件描述符
	fd *os.File
}

// AcquireFileLock 获取锁方法，主要通过syscall中的Flock方法
func AcquireFileLock(path string, readOnly bool) (*FileLockGuard, error) {
	var flag = os.O_RDWR
	if readOnly {
		flag = os.O_RDONLY
	}
	file, err := os.OpenFile(path, flag, 0)
	//错误判断
	if os.IsNotExist(err) {
		file, err = os.OpenFile(path, flag|os.O_CREATE, 0644)
	}
	if err != nil {
		return nil, err
	}

	//如果不是只读就是排他锁
	var how = syscall.LOCK_EX | syscall.LOCK_NB
	//否则就是共享锁
	if readOnly {
		how = syscall.LOCK_SH | syscall.LOCK_NB
	}
	if err := syscall.Flock(int(file.Fd()), how); err != nil {
		return nil, err
	}
	return &FileLockGuard{file}, nil
}

// SyncDir 对文件夹进行同步落盘
func SyncDir(path string) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	err = fd.Sync()
	if err != nil {
		return err
	}
	err = fd.Close()
	if err != nil {
		return err
	}
	return nil
}

// Release 对FileGuard对象的锁进行释放资源
func (fl *FileLockGuard) Release() error {
	how := syscall.LOCK_NB | syscall.LOCK_UN
	err := syscall.Flock(int(fl.fd.Fd()), how)
	if err != nil {
		return err
	}
	return fl.fd.Close()
}
