package kerlong

import (
	"math/rand"
)

//RandInt64 随机一个int64数字
func RandInt64(min, max int64) int64 {
	if min >= max {
		return max
	}
	return rand.Int63n(max-min) + min
}

//RandInt32 随机一个int32数字
func RandInt32(min, max int32) int32 {
	if min >= max {
		return max
	}
	return rand.Int31n(max-min) + min
}

//RandInt16 随机一个int32数字
func RandInt16(min, max int16) int16 {
	if min >= max {
		return max
	}
	return int16(rand.Int31n(int32(max-min)) + int32(min))
}
