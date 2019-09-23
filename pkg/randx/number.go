package randx

import (
	"math/rand"
	"time"
)

func GetInt(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}
