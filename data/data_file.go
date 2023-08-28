package data

import "KVGo/fio"

const DataFileNameSuffix = ".data"

// DataFile 数据文件抽象
type DataFile struct {
	FileId    uint32        // 文件 id
	WriteOff  int64         // 文件写到了那个位置, 文件偏移
	IoManager fio.IOManager // 数据读写接口
}

// OpenDataFile 打开数据文件
func OpenDataFile(dirPath string, fileId uint32) (*DataFile, error) {
	return nil, nil
}

func (df *DataFile) ReadLogRecord(offset int64) (*LogRecord, int64, error) {
	return nil, 0, nil
}

func (df *DataFile) Write(buf []byte) error {
	return nil
}

func (df *DataFile) Sync() error {
	return nil
}
