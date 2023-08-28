package KVGo

import (
	"KVGo/data"
	"KVGo/index"
	"errors"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// DB 存储引擎实例
type DB struct {
	options    Options
	mu         *sync.RWMutex
	fileIds    []int                     // 文件 id, 只能在加载索引的地方使用
	activeFile *data.DataFile            // 当前活跃的数据文件, 可用于写入
	olderFiles map[uint32]*data.DataFile //旧的数据文件, 只读
	index      index.Indexer             // 内存索引
}

// Open 打开存储引擎实例
func Open(options Options) (*DB, error) {
	// 校验配置项
	if err := checkOptions(options); err != nil {
		return nil, err
	}

	// 检查配置项中数据目录是否存在, 如果不存在则创建目录
	if _, err := os.Stat(options.DirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(options.DirPath, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// 初始化 DB 实例
	db := &DB{
		options:    options,
		mu:         new(sync.RWMutex),
		olderFiles: make(map[uint32]*data.DataFile),
		index:      index.NewIndexer(options.IndexType),
	}

	// 加载数据文件
	if err := db.loadDataFiles(); err != nil {
		return nil, err
	}

	// 从数据文件中加载索引到内存
	if err := db.loadIndexFromDataFiles(); err != nil {
		return nil, err
	}

	return db, nil
}

// Put 写入 Key/Value 数据
func (db *DB) Put(key []byte, value []byte) error {
	// 判断 key 是否有效
	if len(key) == 0 {
		return ErrKeyIsEmpty
	}

	// 构造 LogRecord 结构体
	logRecord := &data.LogRecord{
		Key:   key,
		Value: value,
		Type:  data.LogRecordNormal,
	}

	// 追加写入到当前活跃文件
	pos, err := db.appendLogRecord(logRecord)
	if err != nil {
		return err
	}

	// 更新内存索引
	if ok := db.index.Put(key, pos); !ok {
		return ErrIndexUpdateFailed
	}
	return nil
}

func (db *DB) Get(key []byte) ([]byte, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// 判断 key 的有效性
	if len(key) == 0 {
		return nil, ErrKeyIsEmpty
	}

	// 从内存数据结构中取出 key 对应的索引信息
	logRecordPos := db.index.Get(key)
	if logRecordPos == nil {
		return nil, ErrKeyNotFound
	}

	// 根据文件 id 找到对应数据文件
	var dataFile *data.DataFile
	if db.activeFile.FileId == logRecordPos.Fid {
		dataFile = db.activeFile
	} else {
		dataFile = db.olderFiles[logRecordPos.Fid]
	}

	if dataFile == nil {
		return nil, ErrDataFileNotFound
	}

	// 根据偏移量读取数据
	logRecord, _, err := dataFile.ReadLogRecord(logRecordPos.Offset)
	if err != nil {
		return nil, err
	}

	if logRecord.Type == data.LogRecordDeleted {
		return nil, err
	}

	return logRecord.Value, nil
}

// appendLogRecord 追加写数据到活跃文件
func (db *DB) appendLogRecord(logRecord *data.LogRecord) (*data.LogRecordPos, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// 判断当前活跃文件是否存在
	if db.activeFile == nil {
		if err := db.setActiveDataFile(); err != nil {
			return nil, err
		}
	}

	// 写入数据编码
	encRecord, size := data.EncodeLogRecord(logRecord)

	// 如果写入的数据已经达到了活跃文件的阈值, 则关闭活跃文件, 并打开新的文件
	if db.activeFile.WriteOff+size > db.options.DataFileSize {
		// 先持久化数据文件, 保证已有的数据持久到磁盘中
		if err := db.activeFile.Sync(); err != nil {
			return nil, err
		}

		// 将当前活跃文件转换为旧的数据文件
		db.olderFiles[db.activeFile.FileId] = db.activeFile

		// 打开新的数据文件
		if err := db.setActiveDataFile(); err != nil {
			return nil, err
		}
	}

	writeOff := db.activeFile.WriteOff
	if err := db.activeFile.Write(encRecord); err != nil {
		return nil, err
	}

	// 根据用户配置判断是否需要持久化
	if db.options.SyncWrites {
		if err := db.activeFile.Sync(); err != nil {
			return nil, err
		}
	}

	// 构造内存索引信息
	pos := &data.LogRecordPos{Fid: db.activeFile.FileId, Offset: writeOff}
	return pos, nil
}

// setActiveDataFile 设置当前活跃文件, 在访问此方法前必须持有互斥锁
func (db *DB) setActiveDataFile() error {
	var initialFileId uint32 = 0
	if db.activeFile != nil {
		initialFileId = db.activeFile.FileId + 1
	}

	// 打开新的数据文件
	dataFile, err := data.OpenDataFile(db.options.DirPath, initialFileId)
	if err != nil {
		return err
	}

	db.activeFile = dataFile
	return nil
}

// loadDataFiles 从磁盘中加载数据文件
func (db *DB) loadDataFiles() error {
	dirEntries, err := os.ReadDir(db.options.DirPath)
	if err != nil {
		return err
	}

	// 遍历目录中的所有文件, 找到所有以 .data 结尾的文件
	var fileIds []int
	for _, entry := range dirEntries {
		if strings.HasSuffix(entry.Name(), data.DataFileNameSuffix) {
			splitName := strings.Split(entry.Name(), ".")
			fileId, err := strconv.Atoi(splitName[0])
			// 数据目录可能损坏
			if err != nil {
				return ErrDataDirectoryCorrupted
			}
			fileIds = append(fileIds, fileId)
		}
	}

	// 对文件 id 排序, 从大到小依次排序
	sort.Ints(fileIds)
	db.fileIds = fileIds

	// 遍历每个文件 id, 打开对应的数据文件
	for i, fid := range fileIds {
		dataFile, err := data.OpenDataFile(db.options.DirPath, uint32(fid))
		if err != nil {
			return err
		}
		if i == len(fileIds)-1 {
			db.activeFile = dataFile
		} else {
			db.olderFiles[uint32(fid)] = dataFile
		}
	}
	return nil
}

// 从数据文件中加载索引
func (db *DB) loadIndexFromDataFiles() error {
	// 没有文件, 空数据库
	if len(db.fileIds) == 0 {
		return nil
	}

	// 遍历所有的文件 id, 处理文件中的记录
	for i, fid := range db.fileIds {
		var fileId = uint32(fid)
		var dataFile *data.DataFile
		if fileId == db.activeFile.FileId {
			dataFile = db.activeFile
		} else {
			dataFile = db.olderFiles[fileId]
		}

		var offset int64 = 0
		for {
			logRecord, size, err := dataFile.ReadLogRecord(offset)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}

			// 构造内存索引并保存
			logRecordPos := &data.LogRecordPos{Fid: fileId, Offset: offset}
			if logRecord.Type == data.LogRecordDeleted {
				db.index.Delete(logRecord.Key)
			} else {
				db.index.Put(logRecord.Key, logRecordPos)
			}

			// 递增 offset
			offset += size
		}

		if i == len(db.fileIds)-1 {
			db.activeFile.WriteOff = offset
		}
	}
	return nil
}

// checkOptions 检查启动配置项
func checkOptions(options Options) error {
	if options.DirPath == "" {
		return errors.New("database dir path is empty")
	}

	if options.DataFileSize <= 0 {
		return errors.New("database data file size must be greater than 0")
	}

	return nil
}
