package logfile

import (
	"errors"
	"fmt"
	"path/filepath"
	"randb/ioselector"
	"sync"
)

var (
	//ErrInvalidCrc 无效的crc校验
	ErrInvalidCrc = errors.New("logfile: invalid crc")

	// ErrWriteSizeNotEqual 写入大小不等于entry的大小
	ErrWriteSizeNotEqual = errors.New("logfile: write size is not equal to entry size")

	// ErrEndofEntry 日志文件的末尾
	ErrEndofEntry = errors.New("logfile: end of entry in log file")

	// ErrUnsupportedIoType 不支持的IO类型，现在仅仅支持mmap和fileIO
	ErrUnsupportedIoType = errors.New("unsupported io type")

	// ErrUnsupportedLogFileType 不支持的log file类型，现在仅仅支持WAL和ValueLog
	ErrUnsupportedLogFileType = errors.New("unsupported log file type")
)

const (
	//文件id从0开始
	InitialLogFileId = 0

	//日志文件前缀是log.
	FilePrefix = "rance.log."
)


// FileType 文件类型需要对数据结构进行分类
type FileType int8

const (
	Strs FileType = 0
	List FileType = 1
	Hash FileType = 2
	Sets FileType = 3
	Zset FileType = 4
)

var (
	// FileNamesMap 文件类型与名字前缀的映射
	FileNamesMap = map[FileType]string{
		Strs: "rance.log.strs.",
		List: "rance.log.list.",
		Hash: "rance.log.hash.",
		Sets: "rance.log.sets.",
		Zset: "rance.log.zset.",
	}

	// FileTypesMap 文件类型与文件名字符串的映射
	FileTypesMap = map[FileType]string{
		Strs: "strs",
		List: "list",
		Hash: "hash",
		Sets: "sets",
		Zset: "zset",
	}
)

// IOtype 代表了文件IO的不同方式，当前只支持FIleIO和MMap
type IOType int8 

const (
	FileIO IOType = 0
	MMap   IOType = 1
)

// LogFile 是对磁盘文件的抽象，entry的读写都会通过LogFile进行
// 日志文件的记录，由于采用的是日志记录进行具体数据库的更新和删除，很重要
type LogFile struct {
	//读写锁
	sync.RWMutex
	Fid        uint32
	WriteAt    int64
	IoSelector ioselector.IOSelector
}


// OpenLogFile 打开一个存在的日志文件，或者新建一个日志文件
// fileSize必须是正的，会根据ioType来决定用哪种方式打开
func OpenLogFile(path string, fid uint32, fileSize int64, fileType FileType, ioType IOType) (lf *LogFile, err error) {
	lf = &LogFile{Fid: fid}
	fileName, err := lf.getLogFileName(path, fid, fileType)
	if err != nil {
		return nil, err
	}
	var selector ioselector.IOSelector
	switch ioType {
	case FileIO:
		if selector, err = ioselector.NewFileIOSelector(fileName, fileSize); err != nil {
			return nil, err
		}
	case MMap:
		if selector, err = ioselector.NewMMapSelector(fileName, fileSize); err != nil {
			return nil, err
		}
	}
	lf.IoSelector = selector
	return
}

func (lf *LogFile) readLogEntry(offset int64) (*Entry, int64, error) {
	
}

// getLogFileName 获取日志文件名
func (lf *LogFile) getLogFileName(path string, fid uint32, fileType FileType) (name string, err error) {
	prefix, ok := FileNamesMap[fileType]
	if !ok {
		return "", ErrUnsupportedLogFileType
	}
	fileName := prefix + fmt.Sprintf("%09d", fid)
	name = filepath.Join(path, fileName)
	return
}