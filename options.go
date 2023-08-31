package kvgo

import "os"

type Options struct {
	// 数据库数据目录
	DirPath string

	// 数据文件的大小
	DataFileSize int64

	// 累计写入多少字节后进行持久化
	BytesPerSync uint

	// 每次写数据是否持久化
	SyncWrites bool

	// 索引类型
	IndexType IndexerType

	// 是否需要在启动时, 使用 mmap 加载
	MMapAtStartup bool

	// 数据文件合并的阈值
	DataFileMergeRatio float32
}

// IteratorOptions 索引迭代器配置项
type IteratorOptions struct {
	// 遍历前缀为指定值的 Key, 默认为空
	Prefix []byte
	// 是否反向遍历, 默认 false 是正向
	Reverse bool
}

// WriteBatchOptions 批量提交配置项
type WriteBatchOptions struct {
	// 一个 Batch 中最大的数据量
	MaxBatchNum uint

	// 提交时是否 Sync 持久化
	SyncWrites bool
}

type IndexerType = int8

const (
	// BTree 索引
	BTree IndexerType = iota + 1

	// ART 自适应基数树索引
	ART

	// BPlusTree B+树 索引
	BPlusTree
)

var DefaultOptions = Options{
	DirPath:            os.TempDir(),
	DataFileSize:       256 * 1024 * 1024, // 256MB
	SyncWrites:         false,
	BytesPerSync:       0,
	IndexType:          ART,
	MMapAtStartup:      true,
	DataFileMergeRatio: 0.5,
}

var DefaultIteratorOptions = IteratorOptions{
	Prefix:  nil,
	Reverse: false,
}

var DefaultWriteBatchOptions = WriteBatchOptions{
	MaxBatchNum: 10000,
	SyncWrites:  true,
}
