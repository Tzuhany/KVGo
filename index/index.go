package index

import (
	"KVGo/data"
	"bytes"
	"github.com/google/btree"
)

// Indexer 抽象索引接口, 后续接入其他数据结构, 直接实现这个接口即可
type Indexer interface {
	// Put 向索引中存入 key 存储位置信息
	Put(key []byte, pos *data.LogRecordPos) bool

	// Get 根据 key 取出对应索引位置信息
	Get(key []byte) *data.LogRecordPos

	// Delete 根据 key 删除对应索引位置信息
	Delete(key []byte) bool
}

type IndexType = int8

const (
	// Btree 索引
	Btree IndexType = iota + 1

	// ART Adaptive Radix Tree 自适应基数树索引
	ART
)

func NewIndexer(typ IndexType) Indexer {
	switch typ {
	case Btree:
		return NewBtree()
	case ART:
		return nil
	default:
		panic("unsupported index type")
	}
}

type Item struct {
	key []byte
	pos *data.LogRecordPos
}

// Less 自定义 btree 中 key 的比较方法(排序规则)
func (ai *Item) Less(bi btree.Item) bool {
	return bytes.Compare(ai.key, bi.(*Item).key) == -1
}
