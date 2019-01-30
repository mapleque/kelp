package str

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func RandMd5() string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	prefix := strconv.Itoa(rand.Intn(10000))
	surfix := strconv.Itoa(rand.Intn(10000))
	return Md5(fmt.Sprintf("%s%s%s", prefix, timestamp, surfix))
}
