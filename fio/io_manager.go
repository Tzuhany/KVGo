package fio

const DataFilePerm = 0644

// IOManager 抽象 IO 管理接口, 可以接入不同的 IO 类型, 目前支持标准文件 IO
type IOManager interface {
	// Read 从文件给定位置读取数据
	Read([]byte, int64) (int, error)

	// Write 写入字节数组到文件中
	Write([]byte) (int, error)

	// Sync 内存缓冲区的数据持久化到磁盘中
	Sync() error

	// Close 关闭文件
	Close() error
}
