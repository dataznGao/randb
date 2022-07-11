package randb

import "time"

// DataIndexMode 数据索引的格式
type DataIndexMode int

const (
	// KeyValueMemMode 键值对都放在内存中的模式，速度很快
	KeyValueMemMode DataIndexMode = iota

	// KeyOnlyMemMode 只有键在内存中的模式，值需要从日志文件中搜寻出来
	KeyOnlyMemMode
)

// IOType 有两种IO类型可以选: FileIO(标准文件IO) and MMap(内存映射).
type IOType int8

const (
	// FileIO standard file io.
	FileIO IOType = iota
	// MMap Memory Map.
	MMap
)

// Options 启动数据库时的选项
type Options struct {
	// DBPath 数据库存放路径，如果不指定会放在默认路径:TODO
	DBPath string

	// IndexMode 索引模式， 支持 KeyValueMemMode 和 KeyOnlyMemMode.
	// 注意这种模式仅仅影响键值对 String, 而不是 List, Hash, Set 和 ZSet.
	// 默认值是 KeyOnlyMemMode.
	IndexMode DataIndexMode

	// IoType 文件读写的方法, 支持 FileIO 和 MMap.
	// 默认值是 FileIO.
	IoType IOType

	// Sync 是是否完成同步到磁盘的标志
	// 如果是false，那不能保证一定会同步到磁盘
	// 默认值是 false
	Sync bool

	// LogFileGCInterval 一个后台垃圾回收程序会定时进行垃圾回收
	// 它会将满足条件的日志进行删除，然后一条一条进行覆写
	// 默认值是 8小时
	LogFileGCInterval time.Duration

	// LogFileGCRatio if 如果discard日志文件的数量超过了这个比例, 就会被垃圾回收机制挑出来
	// 推荐的比例是0.5
	// 默认值是 0.5.
	LogFileGCRatio float64

	// LogFileSizeThreshold 单个日志文件的大小上限
	// 默认值是 512MB.
	LogFileSizeThreshold int64

	// DiscardBufferSize 当一个key被更新的时候，一条discard数据就会被发送到通道中
	// 这个选项就是调整该通道的
	// 默认值是 8MB.
	DiscardBufferSize int
}

// DefaultOptions 生成打开RanDB数据库的默认选项
func DefaultOptions(path string) Options {
	return Options{
		DBPath:               path,
		IndexMode:            KeyOnlyMemMode,
		IoType:               FileIO,
		Sync:                 false,
		LogFileGCInterval:    time.Hour * 8,
		LogFileGCRatio:       0.5,
		LogFileSizeThreshold: 512 << 20,
		DiscardBufferSize:    8 << 20,
	}
}
