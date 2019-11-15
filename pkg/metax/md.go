package metax

import "errors"

type MetaOption func(md MetaData) MetaData

type MetaData interface {
	Get(interface{}) (interface{}, bool)
	Set(k, v interface{}) error
}

type MetaMap map[interface{}]interface{}

func (mm MetaMap) Get(k interface{}) (interface{}, bool) {
	v, b := mm[k]
	return v, b
}

func (mm MetaMap) Set(k, v interface{}) error {
	mm[k] = v
	return nil
}

var EmptyMD = newMD(false, false)

type MD struct {
	overridable bool
	appendable  bool
	data        map[interface{}]interface{}
}

func (md MD) Get(k interface{}) (interface{}, bool) {
	v, b := md.data[k]
	return v, b
}

func (md MD) Set(k, v interface{}) error {
	if _, ok := md.data[k]; ok {
		if !md.overridable {
			return errors.New("this MD is not overridable")
		}
	} else {
		if !md.appendable {
			return errors.New("this MD is not appendable")
		}
	}
	md.data[k] = v
	return nil
}

func ApplyOption(md MetaData, options ...MetaOption) MetaData {
	if md == nil {
		md = newMD(true, true)
	}
	if len(options) > 0 {
		for _, opt := range options {
			md = opt(md)
		}
	}
	return md
}

func FillMd(md MetaData, data ...map[interface{}]interface{}) {
	if len(data) > 0 {
		for _, d := range data {
			if len(d) > 0 {
				for k, v := range d {
					md.Set(k, v)
				}
			}
		}
	}
}

func newMD(appendable, overridable bool, data ...map[interface{}]interface{}) *MD {
	md := make(map[interface{}]interface{})
	if len(data) > 0 {
		for _, d := range data {
			if len(d) > 0 {
				for k, v := range d {
					md[k] = v
				}
			}
		}
	}
	return &MD{
		overridable: overridable,
		appendable:  appendable,
		data:        md,
	}
}

func ReadonlyMD(data ...map[interface{}]interface{}) *MD {
	return newMD(false, false, data...)
}

func AppendOnlyMD(data ...map[interface{}]interface{}) *MD {
	return newMD(true, false, data...)
}

func NewMD(data ...map[interface{}]interface{}) *MD {
	return newMD(true, true, data...)
}

type WrapMD struct {
	parent MetaData
	md     MetaData
}

func (wm WrapMD) Get(k interface{}) (interface{}, bool) {
	if v, ok := wm.md.Get(k); ok {
		return v, ok
	} else if wm.parent != nil {
		return wm.parent.Get(k)
	} else {
		return v, ok
	}
}

func (wm WrapMD) Set(k, v interface{}) error {
	return wm.md.Set(k, v)
}

func NewWrapMD(parent MetaData, data ...map[interface{}]interface{}) *WrapMD {
	return &WrapMD{
		parent: parent,
		md:     NewMD(data...),
	}
}

func Wrap(parent MetaData, data MetaData) *WrapMD {
	return &WrapMD{
		parent: parent,
		md:     data,
	}
}
