package randx

import (
	"encoding/base64"
	"math/rand"
)

func GetString(s int) (string, error) {
	b, err := GetBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func GetBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

type CharPool struct {
	randPool *Pool
}

func (cp CharPool) Get(length int) string {
	results := make([]rune, length)
	cp.randPool.ResetIndex()
	for i := 0; i < length; i++ {
		randObj := cp.randPool.Get()
		if r, ok := randObj.(rune); ok {
			results[i] = r
		}
	}
	return string(results)
}

func NewCharPool(css ...*CharSource) *CharPool {
	params := make([]interface{}, len(css))
	for i, param := range css {
		params[i] = param
	}
	return &CharPool{randPool: NewPool(params...)}
}

type CharSource struct {
	set []rune
	w   int
}

func (cs CharSource) Get() interface{} {
	if len(cs.set) == 0 {
		panic("CharSource is empty!!!!")
	}
	i := rand.Intn(len(cs.set))
	return cs.set[i]
}

func (cs CharSource) Weight() int {
	return cs.w
}

func NewCharSource(text string, weight int) *CharSource {
	return &CharSource{
		[]rune(text),
		weight,
	}
}
