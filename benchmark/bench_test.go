package benchmark

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"kvgo"
	"kvgo/utils"
	"math/rand"
	"os"
	"testing"
	"time"
)

var db *kvgo.DB

func init() {
	// 初始化用于基准测试的存储引擎
	var err error
	options := kvgo.DefaultOptions
	dir, _ := os.MkdirTemp("", "kvgo-benchmark")
	options.DirPath = dir
	db, err = kvgo.Open(options)
	if err != nil {
		panic(fmt.Sprintf("failed to open db: %v", err))
	}
}

func Benchmark_Put(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := db.Put(utils.GetTestKey(i), utils.RandomValue(1024))
		assert.Nil(b, err)
	}
}

func Benchmark_Get(b *testing.B) {
	for i := 0; i < 100000; i++ {
		err := db.Put(utils.GetTestKey(i), utils.RandomValue(1024))
		assert.Nil(b, err)
	}

	rand.Seed(time.Now().Unix())
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := db.Get(utils.GetTestKey(rand.Int()))
		if err != nil && !errors.Is(err, kvgo.ErrKeyNotFound) {
			b.Fatal(err)
		}
	}
}

func Benchmark_Delete(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err := db.Delete(utils.GetTestKey(rand.Int()))
		assert.Nil(b, err)
	}
}
