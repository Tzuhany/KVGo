package data

type LogRecordType = byte

const (
	LogRecordNormal LogRecordType = iota
	LogRecordDeleted
)

// LogRecord 写入到数据文件的记录
// 因为数据文件中的数据是追加写入的, 类似于日志的格式, 因此叫做日志
type LogRecord struct {
	Key   []byte
	Value []byte
	Type  LogRecordType
}

// LogRecordPos 数据内存索引, 主要描述数据在磁盘上的位置
type LogRecordPos struct {
	Fid    uint32 // 文件 id, 表示数据存储在哪个文件中
	Offset int64  // 偏移量, 表示数据存储在文件中的哪个位置
}

// EncodeLogRecord 对 LogRecord 进行编码, 返回字节数组及长度
func EncodeLogRecord(logRecord *LogRecord) ([]byte, int64) {
	return nil, 0
}
