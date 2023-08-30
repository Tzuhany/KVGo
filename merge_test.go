package kvgo

import (
	"github.com/stretchr/testify/assert"
	"kvgo/utils"
	"testing"
)

func TestDB_Merge(t *testing.T) {
	opts := DefaultOptions
	dir := "/tmp/kvgo-merge-1"
	opts.DirPath = dir
	opts.DataFileSize = 64 * 1024 * 1024
	db, err := Open(opts)
	//defer destroyDB(db)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	//keys := db.ListKeys()
	//t.Log(len(keys))

	val, err := db.Get(utils.GetTestKey(99033))
	t.Log(string(val))
	t.Log(err)

	//for i := 0; i < 500000; i++ {
	//	db.Put(utils.GetTestKey(i), utils.RandomValue(128))
	//}
	//for i := 0; i < 500000; i++ {
	//	if i == 99033 {
	//		db.Put(utils.GetTestKey(i), utils.RandomValue(128))
	//	} else {
	//		db.Delete(utils.GetTestKey(i))
	//	}
	//}
	//db.Merge()
}
