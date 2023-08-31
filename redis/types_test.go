package redis

import (
	"github.com/stretchr/testify/assert"
	"kvgo"
	"kvgo/utils"
	"os"
	"testing"
	"time"
)

func TestRedisDataStructure_Get(t *testing.T) {
	opts := kvgo.DefaultOptions
	dir, _ := os.MkdirTemp("", "kvgo-redis-get")
	opts.DirPath = dir
	rds, err := NewRedisDataStructure(opts)
	assert.Nil(t, err)

	err = rds.Set(utils.GetTestKey(1), 0, utils.RandomValue(100))
	assert.Nil(t, err)
	err = rds.Set(utils.GetTestKey(2), time.Second*5, utils.RandomValue(100))
	assert.Nil(t, err)

	val1, err := rds.Get(utils.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, val1)

	val2, err := rds.Get(utils.GetTestKey(2))
	assert.Nil(t, err)
	assert.NotNil(t, val2)

	_, err = rds.Get(utils.GetTestKey(33))
	assert.Equal(t, kvgo.ErrKeyNotFound, err)
}

func TestRedisDataStructure_Del_Type(t *testing.T) {
	opts := kvgo.DefaultOptions
	dir, _ := os.MkdirTemp("", "kvgo-redis-del-type")
	opts.DirPath = dir
	rds, err := NewRedisDataStructure(opts)
	assert.Nil(t, err)

	// del
	err = rds.Del(utils.GetTestKey(11))
	assert.Nil(t, err)

	err = rds.Set(utils.GetTestKey(1), 0, utils.RandomValue(100))
	assert.Nil(t, err)

	// type
	typ, err := rds.Type(utils.GetTestKey(1))
	assert.Nil(t, err)
	assert.Equal(t, String, typ)

	err = rds.Del(utils.GetTestKey(1))
	assert.Nil(t, err)

	_, err = rds.Get(utils.GetTestKey(1))
	assert.Equal(t, kvgo.ErrKeyNotFound, err)
}

func TestRedisDataStructure_HGet(t *testing.T) {
	opts := kvgo.DefaultOptions
	dir, _ := os.MkdirTemp("", "kvgo-redis-hget")
	opts.DirPath = dir
	rds, err := NewRedisDataStructure(opts)
	assert.Nil(t, err)

	ok1, err := rds.HSet(utils.GetTestKey(1), []byte("field1"), utils.RandomValue(100))
	assert.Nil(t, err)
	assert.True(t, ok1)

	v1 := utils.RandomValue(100)
	ok2, err := rds.HSet(utils.GetTestKey(1), []byte("field1"), v1)
	assert.Nil(t, err)
	assert.False(t, ok2)

	v2 := utils.RandomValue(100)
	ok3, err := rds.HSet(utils.GetTestKey(1), []byte("field2"), v2)
	assert.Nil(t, err)
	assert.True(t, ok3)

	val1, err := rds.HGet(utils.GetTestKey(1), []byte("field1"))
	assert.Nil(t, err)
	assert.Equal(t, v1, val1)

	val2, err := rds.HGet(utils.GetTestKey(1), []byte("field2"))
	assert.Nil(t, err)
	assert.Equal(t, v2, val2)

	_, err = rds.HGet(utils.GetTestKey(1), []byte("field-not-exist"))
	assert.Equal(t, kvgo.ErrKeyNotFound, err)
}

func TestRedisDataStructure_HDel(t *testing.T) {
	opts := kvgo.DefaultOptions
	dir, _ := os.MkdirTemp("", "kvgo-redis-hdel")
	opts.DirPath = dir
	rds, err := NewRedisDataStructure(opts)
	assert.Nil(t, err)

	del1, err := rds.HDel(utils.GetTestKey(200), nil)
	assert.Nil(t, err)
	assert.False(t, del1)

	ok1, err := rds.HSet(utils.GetTestKey(1), []byte("field1"), utils.RandomValue(100))
	assert.Nil(t, err)
	assert.True(t, ok1)

	v1 := utils.RandomValue(100)
	ok2, err := rds.HSet(utils.GetTestKey(1), []byte("field1"), v1)
	assert.Nil(t, err)
	assert.False(t, ok2)

	v2 := utils.RandomValue(100)
	ok3, err := rds.HSet(utils.GetTestKey(1), []byte("field2"), v2)
	assert.Nil(t, err)
	assert.True(t, ok3)

	del2, err := rds.HDel(utils.GetTestKey(1), []byte("field1"))
	assert.Nil(t, err)
	assert.True(t, del2)
}

func TestRedisDataStructure_SIsMember(t *testing.T) {
	opts := kvgo.DefaultOptions
	dir, _ := os.MkdirTemp("", "kvgo-redis-sismember")
	opts.DirPath = dir
	rds, err := NewRedisDataStructure(opts)
	assert.Nil(t, err)

	ok, err := rds.SAdd(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.True(t, ok)
	ok, err = rds.SAdd(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.False(t, ok)
	ok, err = rds.SAdd(utils.GetTestKey(1), []byte("val-2"))
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = rds.SIsMember(utils.GetTestKey(2), []byte("val-1"))
	assert.Nil(t, err)
	assert.False(t, ok)
	ok, err = rds.SIsMember(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.True(t, ok)
	ok, err = rds.SIsMember(utils.GetTestKey(1), []byte("val-2"))
	assert.Nil(t, err)
	assert.True(t, ok)
	ok, err = rds.SIsMember(utils.GetTestKey(1), []byte("val-not-exist"))
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestRedisDataStructure_SRem(t *testing.T) {
	opts := kvgo.DefaultOptions
	dir, _ := os.MkdirTemp("", "kvgo-go-redis-srem")
	opts.DirPath = dir
	rds, err := NewRedisDataStructure(opts)
	assert.Nil(t, err)

	ok, err := rds.SAdd(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.True(t, ok)
	ok, err = rds.SAdd(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.False(t, ok)
	ok, err = rds.SAdd(utils.GetTestKey(1), []byte("val-2"))
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = rds.SRem(utils.GetTestKey(2), []byte("val-1"))
	assert.Nil(t, err)
	assert.False(t, ok)
	ok, err = rds.SRem(utils.GetTestKey(1), []byte("val-2"))
	assert.Nil(t, err)
	assert.True(t, ok)

	ok, err = rds.SIsMember(utils.GetTestKey(1), []byte("val-2"))
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestRedisDataStructure_LPop(t *testing.T) {
	opts := kvgo.DefaultOptions
	dir, _ := os.MkdirTemp("", "kvgo-redis-lpop")
	opts.DirPath = dir
	rds, err := NewRedisDataStructure(opts)
	assert.Nil(t, err)

	res, err := rds.LPush(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), res)
	res, err = rds.LPush(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.Equal(t, uint32(2), res)
	res, err = rds.LPush(utils.GetTestKey(1), []byte("val-2"))
	assert.Nil(t, err)
	assert.Equal(t, uint32(3), res)

	val, err := rds.LPop(utils.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, val)
	val, err = rds.LPop(utils.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, val)
	val, err = rds.LPop(utils.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, val)
}

func TestRedisDataStructure_RPop(t *testing.T) {
	opts := kvgo.DefaultOptions
	dir, _ := os.MkdirTemp("", "kvgo-redis-rpop")
	opts.DirPath = dir
	rds, err := NewRedisDataStructure(opts)
	assert.Nil(t, err)

	res, err := rds.RPush(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), res)
	res, err = rds.RPush(utils.GetTestKey(1), []byte("val-1"))
	assert.Nil(t, err)
	assert.Equal(t, uint32(2), res)
	res, err = rds.RPush(utils.GetTestKey(1), []byte("val-2"))
	assert.Nil(t, err)
	assert.Equal(t, uint32(3), res)

	val, err := rds.RPop(utils.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, val)
	val, err = rds.RPop(utils.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, val)
	val, err = rds.RPop(utils.GetTestKey(1))
	assert.Nil(t, err)
	assert.NotNil(t, val)
}

func TestRedisDataStructure_ZSet(t *testing.T) {
	opts := kvgo.DefaultOptions
	dir := "/tmp/kvgo-redis-zset"
	opts.DirPath = dir
	rds, err := NewRedisDataStructure(opts)
	assert.Nil(t, err)

	res, err := rds.ZAdd(utils.GetTestKey(1), 65, []byte("val-1"))
	t.Log(res, err)
	res, err = rds.ZAdd(utils.GetTestKey(1), 88, []byte("val-1"))
	t.Log(res, err)
	res, err = rds.ZAdd(utils.GetTestKey(1), 11, []byte("val-2"))
	t.Log(res, err)
	res, err = rds.ZAdd(utils.GetTestKey(1), 22, []byte("val-5"))
	t.Log(res, err)
	res, err = rds.ZAdd(utils.GetTestKey(1), 11, []byte("val-7"))
	t.Log(res, err)
	res, err = rds.ZAdd(utils.GetTestKey(1), 65, []byte("val-3"))
	t.Log(res, err)
	res, err = rds.ZAdd(utils.GetTestKey(1), 18, []byte("val-9"))
	t.Log(res, err)

	score, err := rds.ZScore(utils.GetTestKey(1), []byte("val-1"))
	t.Log(score, err)
}
