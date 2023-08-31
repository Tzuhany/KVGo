package utils

import (
	"strconv"
)

func Float64FromBytes(bytes []byte) float64 {
	val, _ := strconv.ParseFloat(string(bytes), 64)
	return val
}

func Float64ToBytes(float float64) []byte {
	return []byte(strconv.FormatFloat(float, 'f', -1, 64))
}
