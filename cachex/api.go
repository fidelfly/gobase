package cachex

import "github.com/fidelfly/gox/cachex/bcache"

//export
func NewBuntCache(filename string) *bcache.BuntCache {
	cache, err := bcache.NewCache(filename)
	if err != nil {
		return nil
	}
	return cache
}
