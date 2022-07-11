package randb

import (
	"errors"
	"randb/ioselector"
	"sync"
)

const (
	discardRecordSize = 12
	// 8kb, contains mostly 682 records in file.
	discardFileSize int64 = 2 << 12
	discardFileName       = "discard"
)

// ErrDiscardNoSpace 没有剩余空间存放discard文件
var ErrDiscardNoSpace = errors.New("not enough space can be allocated for the discard file")

type discard struct {
	// 临界锁
	sync.Mutex
	// 仅进入一次
	once *sync.Once
	// 通道，将需要删除的下标点吸收进来
	valChan chan *indexNode
	//文件通道，用来操作文件
	file ioselector.IOSelector
	//释放的文件，记录文件偏移量
	freeList []int64
	//每一个fid对应的偏移量
	location map[uint32]int64
}
