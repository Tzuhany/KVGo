package KVGo

type Options struct {
	// 数据库数据目录
	DirPath string

	// 数据文件的大小
	DataFileSize int64

	// 是否每次写入持久化
	SyncWrites bool

	// 索引类型
	IndexType IndexType
}

type IndexType = int8

const (
	// Btree 索引
	Btree IndexType = iota + 1

	// ART Adaptive Radix Tree 自适应基数树索引
	ART
)
