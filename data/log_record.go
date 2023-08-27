package data

// LogRecordPos 数据内存索引，主要描述数据在磁盘上的位置
type LogRecordPos struct {
	Fid    uint32 // 文件 id，表示数据存储在哪个文件中
	Offset int64  // 偏移量，表示数据存储在文件中的哪个位置
}
