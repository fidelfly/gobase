package bcache

import (
	"strings"

	"github.com/tidwall/buntdb"
)

//export
func NewCache(filename string, portions ...Portion) (*BuntCache, error) {
	db, err := buntdb.Open(filename)
	if err != nil {
		return nil, err
	}
	bc := &BuntCache{db: db}
	bc.AddPortion(DefaultPortion)
	if len(portions) > 0 {
		bc.AddPortion(portions...)
	}
	return bc, nil
}

//export
func NewCacheWithConverter(filename string, converter Converter) (*BuntCache, error) {
	db, err := buntdb.Open(filename)
	if err != nil {
		return nil, err
	}
	bc := &BuntCache{db: db}
	bc.AddPortion(NewDataset("*", converter))
	return bc, nil
}

//export
func NewKey(keys ...string) string {
	return strings.Join(keys, ":")
}
