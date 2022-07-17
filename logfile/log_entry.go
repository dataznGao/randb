package logfile

import (
	"encoding/binary"
	"hash/crc32"
)

// 真正落盘时，是entryHeader  + logEntry这样反复迭代
// logEntry是真正落盘的数据条，entryHeader是元数据头
// 这个元数据头不希望被外界读到
const MaxHeaderSize = 25

// EntryType entry的类型
type EntryType byte

const (
	// TypeDelete 代表这个entry已经被删除了
	TypeDelete EntryType = 1
	// TypeListMeta 代表这个是ListMeta
	TypeListMeta EntryType = 2
)

type EntryBody struct {
	Key       []byte
	Value     []byte
	expiredAt int64
	Type      EntryType
}

type entryHeader struct {
	crc32     uint32
	typ       EntryType
	kSize     uint32
	vSize     uint32
	expiredAt uint32
}

// 存储的时候将整条数据：entryHeader + entryBody都encode成byte
// 读取的时候，只要读取到头部，就能知道对应的key value是啥了
// +-------+--------+----------+------------+-----------+-------+---------+
// |  crc  |  type  | key size | value size | expiredAt |  key  |  value  |
// +-------+--------+----------+------------+-----------+-------+---------+
// |------------------------HEADER----------------------|
//         |--------------------------crc check---------------------------|
func EncodeEntry(e *EntryBody) ([]byte, int) {
	if e == nil {
		return nil, 0
	}

	// 编码头
	header := make([]byte, MaxHeaderSize)
	// crc 32位, 4个字节
	// type
	header[4] = byte(e.Type)
	// key size
	index := 5
	index += binary.PutVarint(header[index:], int64(len(e.Key)))
	// Value size
	index += binary.PutVarint(header[index:], int64(len(e.Value)))
	// expiredAt
	index += binary.PutVarint(header[index:], int64(e.expiredAt))

	var size = index + len(e.Key) + len(e.Value)
	buf := make([]byte, size)
	copy(buf[:index], header[:])
	copy(buf[index:], e.Key)
	copy(buf[index+len(e.Key):], e.Value)

	// crc校验
	crc := crc32.ChecksumIEEE(buf[4:])
	binary.LittleEndian.PutUint32(buf[:4], crc)
	return buf, size
}
