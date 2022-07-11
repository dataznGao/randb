package randb

type DataType = int8

//支持5种不同的数据类型，支持String, List, Hash, Set 和 Sorted Set
const (
	String DataType = 0
	List   DataType = 1
	Hash   DataType = 2
	Set    DataType = 3
	ZSet   DataType = 4
)
