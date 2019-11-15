package bcache

import (
	"encoding/json"

	"github.com/tidwall/buntdb"
)

//https://github.com/tidwall/buntdb

type BuntCache struct {
	db       *buntdb.DB
	portions []Portion
}

type Portion interface {
	Converter
	Match(key string) bool
}

type Converter interface {
	Encode(obj interface{}) (string, error)
	Decode(val string, obj interface{}) error
}

func (bc *BuntCache) findPortion(key string) Portion {
	lop := len(bc.portions)
	if lop > 0 {
		for i := lop - 1; i >= 0; i-- {
			ds := bc.portions[i]
			if ds.Match(key) {
				return ds
			}
		}
	}
	return nil
}

type Dataset struct {
	pattern string
	coder   Converter
	//Indexes map[string]func(a string, b string) bool
}

func (ds *Dataset) Encode(obj interface{}) (string, error) {
	return ds.coder.Encode(obj)
}
func (ds *Dataset) Decode(val string, obj interface{}) error {
	return ds.coder.Decode(val, obj)
}
func (ds *Dataset) Match(key string) bool {
	return buntdb.Match(key, ds.pattern)
}

func NewDataset(p string, coder Converter) *Dataset {
	return &Dataset{p, coder}
}

var JsonConverter = &jsonConverter{}
var DefaultPortion = &Dataset{"*", JsonConverter}

type jsonConverter struct {
}

func (jc *jsonConverter) Encode(obj interface{}) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
func (jc *jsonConverter) Decode(val string, obj interface{}) error {
	return json.Unmarshal([]byte(val), obj)
}

func (bc *BuntCache) SetDefaultConverter(converter Converter) {
	if len(bc.portions) > 0 {
		bc.portions[0] = NewDataset("*", converter)
	} else {
		bc.AddPortion(NewDataset("*", converter))
	}
}

func (bc *BuntCache) AddPortion(portions ...Portion) {
	bc.portions = append(bc.portions, portions...)
}

func (bc *BuntCache) GetDB() *buntdb.DB {
	return bc.db
}

func (bc *BuntCache) Set(key string, obj interface{}) error {
	return bc.db.Update(func(tx *buntdb.Tx) error {
		dp := bc.findPortion(key)
		if dp == nil {
			dp = DefaultPortion
		}
		data, err := dp.Encode(obj)
		_, _, err = tx.Set(key, data, nil)
		return err
	})
}

func (bc *BuntCache) Get(key string) (val string, err error) {
	_ = bc.db.View(func(tx *buntdb.Tx) error {
		val, err = tx.Get(key)
		return nil
	})
	return
}

func (bc *BuntCache) GetObject(key string, val interface{}) (err error) {
	data, err := bc.Get(key)
	if err == nil {
		dp := bc.findPortion(key)
		if dp == nil {
			dp = DefaultPortion
		}
		err = dp.Decode(data, val)
	}
	return
}

func (bc *BuntCache) Delete(key string) error {
	return bc.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(key)
		return err
	})
}

func (bc *BuntCache) Iterate(iterator func(key, val string) bool) {
	_ = bc.db.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys("", iterator)
	})
}

func (bc *BuntCache) CreateJSONIndex(name string, pattern string, paths ...string) {
	if len(paths) == 0 {
		return
	}
	_ = bc.db.Update(func(tx *buntdb.Tx) error {
		indexes := make([]func(string, string) bool, len(paths))
		for i, path := range paths {
			indexes[i] = buntdb.IndexJSON(path)
		}
		return tx.CreateIndex(name, pattern, indexes...)
	})

	return
}
