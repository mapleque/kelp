package util

import (
	"math/rand"
	"strconv"
)

func RandMD5() string {
	seed := rand.Intn(10000)
	randStr := strconv.Itoa(seed)
	return MD5(randStr)
}
