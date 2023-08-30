package kvgo

import (
	"github.com/stretchr/testify/assert"
	"kvgo/utils"
	"os"
	"testing"
)

// 写完未提交

// 写完提交

// 提交后重启，在提交

func TestDB_WriteBatch(t *testing.T) {
	opts := DefaultOptions
	dir, _ := os.MkdirTemp("", "kvgo-wb")
	//dir := "/tmp/batch-test"
	opts.DirPath = dir
	db, err := Open(opts)
	defer destroyDB(db)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	wb := db.NewWriteBatch(DefaultWriteBatchOptions)
	err = wb.Put(utils.GetTestKey(12), utils.RandomValue(10))
	assert.Nil(t, err)
	err = wb.Put(utils.GetTestKey(112), utils.RandomValue(10))
	assert.Nil(t, err)
	err = wb.Put(utils.GetTestKey(121), utils.RandomValue(10))
	assert.Nil(t, err)
	//err = wb.Delete(utils.GetTestKey(12))
	//assert.Nil(t, err)

	_, err = db.Get(utils.GetTestKey(12))
	assert.Equal(t, ErrKeyNotFound, err)

	err = wb.Put(utils.GetTestKey(15), utils.RandomValue(10))
	assert.Nil(t, err)
	err = wb.Commit()
	assert.Nil(t, err)

	val, err := db.Get(utils.GetTestKey(15))
	t.Log(string(val))
	t.Log(err)
}

func TestDB_WriteBatch2(t *testing.T) {
	opts := DefaultOptions
	dir, _ := os.MkdirTemp("", "kvgo-batch-2")
	opts.DirPath = dir
	db, err := Open(opts)
	defer destroyDB(db)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	err = db.Put(utils.GetTestKey(1), utils.RandomValue(10))
	assert.Nil(t, err)

	wb := db.NewWriteBatch(DefaultWriteBatchOptions)
	err = wb.Put(utils.GetTestKey(2), utils.RandomValue(10))
	assert.Nil(t, err)
	err = wb.Delete(utils.GetTestKey(1))
	assert.Nil(t, err)

	err = wb.Commit()
	assert.Nil(t, err)

	err = wb.Put(utils.GetTestKey(11), utils.RandomValue(10))
	assert.Nil(t, err)
	err = wb.Commit()
	assert.Nil(t, err)

	// 重启
	err = db.Close()
	assert.Nil(t, err)

	db2, err := Open(opts)
	assert.Nil(t, err)

	_, err = db2.Get(utils.GetTestKey(1))
	assert.Equal(t, ErrKeyNotFound, err)

	// 校验序列号
	assert.Equal(t, uint64(2), db.seqNo)
}

func TestDB_WriteBatch4(t *testing.T) {
	opts := DefaultOptions
	//dir, _ := os.MkdirTemp("", "kvgo-wb")
	dir := "/tmp/batch-test-3"
	opts.DirPath = dir
	db, err := Open(opts)
	//defer destroyDB(db)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	//wb := db.NewWriteBatch(DefaultWriteBatchOptions)
	//for i := 0; i < 7000; i++ {
	//	wb.Put(utils.GetTestKey(i), utils.RandomValue(40960))
	//}
	//
	//err = wb.Commit()
	//t.Log(err)

	keys := db.ListKeys()
	t.Log(len(keys))
	t.Log(db.seqNo)
}
