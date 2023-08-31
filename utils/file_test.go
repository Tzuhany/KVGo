package utils

import (
	"math"
	"path/filepath"
	"testing"
)

func TestDirSize(t *testing.T) {
	size, err := DirSize(filepath.Join("/tmp/bitcask-go-stat"))
	t.Log(size)
	t.Log(err)

	t.Log(math.MaxInt64 / 1024 / 1024)
	t.Log(math.MaxUint32 / 1024 / 1024)

	diskSize, _ := AvailableDiskSize()
	t.Log(diskSize / 1024 / 1024 / 1024)

	a := 9999999
	b := 9900293

	t.Log(float32(b) / float32(a))
	t.Log(math.MaxFloat32 > math.MaxUint64)
}
