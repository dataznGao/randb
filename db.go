package randb

import (
	"randb/ds/art"
	"randb/ds/zset"
	"randb/flock"
	"randb/logfile"
	"randb/util"
	"sync"

)

type (
	//RanDB 是一个数据库实例
	RanDB struct {
		//activeLogFiles 活跃的日志文件，用map组织
		activeLogFiles map[DataType]*logfile.LogFile
		//archivedLogFiles 已经存档的日志文件，用map的map组织
		archivedLogFiles map[DataType]archivedFiles
		//fidMap 一个日志的索引对应一个int32的数组，这个数组只有在启动的时候需要用到，即便日志文件改变了也不更新这个数组
		fidMap map[DataType][]uint32
		//discards 丢弃的日志文件，也就是需要做删除操作的文件
		discards map[DataType]*discard
		//opts 启动时的选项
		opts Options
		//strIndex 存string的索引树
		strIndex *strIndex
		//listIndex 存list的索引树
		listIndex *listIndex
		//hashIndex 存hash的索引树
		hashIndex *hashIndex
		//hashIndex 存set的索引树
		setIndex *setIndex
		//hashIndex 存zset的索引树
		zsetIndex *zsetIndex
		//mutex 读写锁
		mutex sync.RWMutex
		//fileLock 文件锁
		fileLock *flock.FileLockGuard
		closed   uint32
		gcState  int32
	}

	// 已经存档的日志文件
	archivedFiles map[uint32]*logfile.LogFile

	//文件的索引结点
	indexNode struct {
		value     []byte
		fid       uint32
		offset    int64
		extrySize int
		//到期时间
		expiredAt int64
	}

	listIndex struct {
		//读写锁
		mutex *sync.RWMutex
		//TODO:可变基数树
		trees map[string]*art.AdaptiveRadixTree
	}

	strIndex struct {
		mutex *sync.RWMutex
		tree *art.AdaptiveRadixTree
	}

	hashIndex struct {
		mutex *sync.RWMutex
		trees map[string]*art.AdaptiveRadixTree
	}

	setIndex struct {
		mutex *sync.RWMutex
		murhash *util.Murmur128
		trees map[string]*art.AdaptiveRadixTree
	}

	zsetIndex struct {
		mutex *sync.RWMutex
		indexes *zset.SortedSet
		murhash *util.Murmur128
		trees map[string]*art.AdaptiveRadixTree
	}
)

func newStrsIndex() *strIndex {
	return &strIndex{tree: art.NewART(), mutex: new(sync.RWMutex)}
}

func newListIndex() *listIndex {
	return &listIndex{
		trees: make(map[string]*art.AdaptiveRadixTree), 
		mutex: new(sync.RWMutex),
	}
}

func newHashIndex() *hashIndex {
	return &hashIndex{
		trees: make(map[string]*art.AdaptiveRadixTree), 
		mutex: new(sync.RWMutex),
	}
}

func newSetIndex() *setIndex {
	return &setIndex{
		trees: make(map[string]*art.AdaptiveRadixTree), 
		murhash: util.NewMurmur128(), 
		mutex: new(sync.RWMutex),
	}
}

func newZsetIndex() *zsetIndex {
	return &zsetIndex{
		trees: make(map[string]*art.AdaptiveRadixTree), 
		murhash: util.NewMurmur128(), 
		mutex: new(sync.RWMutex),
	}
}