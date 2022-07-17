package zset

const (
	maxLevel    = 32
	probability = 0.25
)

type EncodeKey func(key, subKey []byte) []byte

type (
	// SortedSet 有序集合的结构
	SortedSet struct {
		record map[string]*SortedSetNode
	}

	// SortedSetNode 有序集合的结点
	SortedSetNode struct {
		dict map[string]*sklNode
		skl  *skipList
	}

	// sklNode 跳表的结点
	sklNode struct {
		member   string
		score    float64
		backward *sklNode
		level    []*sklLevel
	}

	// sklNode 靠它来记录跳几个
	sklLevel struct {
		forward *sklNode
		span    uint64
	}

	// skipList 链式结构
	skipList struct {
		head   *sklNode
		tail   *sklNode
		length int64
		level  int16
	}
)

// 创建一个新的有序表
func New() *SortedSet {
	return &SortedSet{
		make(map[string]*SortedSetNode),
	}
}
