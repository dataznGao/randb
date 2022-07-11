package randb

import (
	"randb/ioselector"
	"sync"
)

// LogFile 日志文件的记录，由于采用的是日志记录进行具体数据库的更新和删除，很重要
type LogFile struct {
	//读写锁
	sync.RWMutex
	Fid        uint32
	WriteAt    int64
	IoSelector ioselector.IOSelector
}
