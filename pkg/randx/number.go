package randx

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//export
func GetInt(n int) int {
	return rand.Intn(n)
}
