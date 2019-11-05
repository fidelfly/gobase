package jmap

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"sync"
)

type JSONMap map[string]interface{}
type Decorator func(string, JSONMap)

func Marshal(v interface{}, recursive bool, decorators ...Decorator) interface{} {
	return convert(reflect.ValueOf(v), convertOpts{recursive: recursive, decorators: decorators})
}

type convertOpts struct {
	recursive  bool
	decorators []Decorator
	path       string
}

type jmConverter func(v reflect.Value, opts convertOpts) interface{}

var converterCache sync.Map

var fieldCache sync.Map // map[reflect.Type][]field

type field struct {
	name      string
	tag       bool
	index     []int
	typ       reflect.Type
	omitEmpty bool

	converter jmConverter
}

// cachedTypeFields is like typeFields but uses a cache to avoid repeated work.
func cachedTypeFields(t reflect.Type) []field {
	if f, ok := fieldCache.Load(t); ok {
		return f.([]field)
	}
	f, _ := fieldCache.LoadOrStore(t, typeFields(t))
	return f.([]field)
}

func typeFields(t reflect.Type) []field {
	// Anonymous fields to explore at the current level and the next.
	var current []field
	next := []field{{typ: t}}

	// Count of queued names for current level and the next.
	var count map[reflect.Type]int
	nextCount := map[reflect.Type]int{}

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool{}

	// Fields found.
	var fields []field

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true

			// Scan f.typ for fields to include.
			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				isUnexported := sf.PkgPath != ""
				if sf.Anonymous {
					t0 := sf.Type
					if t0.Kind() == reflect.Ptr {
						t0 = t0.Elem()
					}
					if isUnexported && t0.Kind() != reflect.Struct {
						// Ignore embedded fields of unexported non-struct types.
						continue
					}
					// Do not ignore embedded fields of unexported struct types
					// since they may have exported fields.
				} else if isUnexported {
					// Ignore unexported non-embedded fields.
					continue
				}
				tag := sf.Tag.Get("json")
				if tag == "-" {
					continue
				}
				name, opts := parseTag(tag)
				if !isValidTag(name) {
					name = ""
				}
				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Ptr {
					// Follow pointer.
					ft = ft.Elem()
				}

				if name != "" || !sf.Anonymous || ft.Kind() != reflect.Struct {
					tagged := name != ""
					if name == "" {
						name = sf.Name
					}
					fieldItem := field{
						name:      name,
						tag:       tagged,
						index:     index,
						typ:       ft,
						omitEmpty: opts.Contains("omitempty"),
					}

					fields = append(fields, fieldItem)
					if count[f.typ] > 1 {
						// If there were multiple instances, add a second,
						// so that the annihilation code will see a duplicate.
						// It only cares about the distinction between 1 or 2,
						// so don't bother generating any more copies.
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}

				// Record new anonymous struct to explore in next round.
				nextCount[ft]++
				if nextCount[ft] == 1 {
					next = append(next, field{name: ft.Name(), index: index, typ: ft})
				}
			}
		}
	}

	sort.Slice(fields, func(i, j int) bool {
		x := fields
		// sort field by name, breaking ties with depth, then
		// breaking ties with "name came from json tag", then
		// breaking ties with index sequence.
		if x[i].name != x[j].name {
			return x[i].name < x[j].name
		}
		if len(x[i].index) != len(x[j].index) {
			return len(x[i].index) < len(x[j].index)
		}
		if x[i].tag != x[j].tag {
			return x[i].tag
		}
		return byIndex(x).Less(i, j)
	})

	// Delete all fields that are hidden by the Go rules for embedded fields,
	// except that fields with JSON tags are promoted.

	// The fields are sorted in primary order of name, secondary order
	// of field index length. Loop over names; for each name, delete
	// hidden fields by choosing the one dominant field that survives.
	out := fields[:0]
	for advance, i := 0, 0; i < len(fields); i += advance {
		// One iteration per name.
		// Find the sequence of fields with the name of this first field.
		fi := fields[i]
		name := fi.name
		for advance = 1; i+advance < len(fields); advance++ {
			fj := fields[i+advance]
			if fj.name != name {
				break
			}
		}
		if advance == 1 { // Only one field with this name
			out = append(out, fi)
			continue
		}
		dominant, ok := dominantField(fields[i : i+advance])
		if ok {
			out = append(out, dominant)
		}
	}

	fields = out
	sort.Sort(byIndex(fields))

	for i := range fields {
		f := &fields[i]
		f.converter = typeConverter(typeByIndex(t, f.index))
	}
	return fields
}

func typeByIndex(t reflect.Type, index []int) reflect.Type {
	for _, i := range index {
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		t = t.Field(i).Type
	}
	return t
}

type byIndex []field

func (x byIndex) Len() int { return len(x) }

func (x byIndex) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x byIndex) Less(i, j int) bool {
	for k, xik := range x[i].index {
		if k >= len(x[j].index) {
			return false
		}
		if xik != x[j].index[k] {
			return xik < x[j].index[k]
		}
	}
	return len(x[i].index) < len(x[j].index)
}

func dominantField(fields []field) (field, bool) {
	// The fields are sorted in increasing index-length order, then by presence of tag.
	// That means that the first field is the dominant one. We need only check
	// for error cases: two fields at top level, either both tagged or neither tagged.
	if len(fields) > 1 && len(fields[0].index) == len(fields[1].index) && fields[0].tag == fields[1].tag {
		return field{}, false
	}
	return fields[0], true
}

func convert(v reflect.Value, opts convertOpts) interface{} {
	return valueConverter(v)(v, opts)
}

func valueConverter(v reflect.Value) jmConverter {
	if !v.IsValid() {
		return nilConverter
	}
	return typeConverter(v.Type())
}

func typeConverter(t reflect.Type) jmConverter {
	if fi, ok := converterCache.Load(t); ok {
		return fi.(jmConverter)
	}
	var (
		wg sync.WaitGroup
		f  jmConverter
	)
	wg.Add(1)
	fi, loaded := converterCache.LoadOrStore(t, jmConverter(func(v reflect.Value, opts convertOpts) interface{} {
		wg.Wait()
		return f(v, opts)
	}))

	if loaded {
		return fi.(jmConverter)
	}

	f = newTypeConverter(t, true)
	wg.Done()
	converterCache.Store(t, f)
	return f
}

var (
	marshalerType     = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

func newCondAddrEncoder(canAddrEnc, elseEnc jmConverter) jmConverter {
	enc := condAddrConverter{canAddrEnc: canAddrEnc, elseEnc: elseEnc}
	return enc.convert
}

type condAddrConverter struct {
	canAddrEnc, elseEnc jmConverter
}

func (ce condAddrConverter) convert(v reflect.Value, opts convertOpts) interface{} {
	if v.CanAddr() {
		return ce.canAddrEnc(v, opts)
	} else {
		return ce.elseEnc(v, opts)
	}
}

func newTypeConverter(t reflect.Type, allowAddr bool) jmConverter {
	if t.Implements(marshalerType) {
		return skipConverter
	}
	if t.Kind() != reflect.Ptr && allowAddr {
		if reflect.PtrTo(t).Implements(marshalerType) {
			return newCondAddrEncoder(skipConverter, newTypeConverter(t, false))
		}
	}

	if t.Implements(textMarshalerType) {
		return skipConverter
	}
	if t.Kind() != reflect.Ptr {
		if reflect.PtrTo(t).Implements(textMarshalerType) && allowAddr {
			return newCondAddrEncoder(skipConverter, newTypeConverter(t, false))
		}
	}
	switch t.Kind() {
	case reflect.Interface:
		return interfaceConverter
	case reflect.Struct:
		return newStructConverter(t)
	case reflect.Map:
		return mapConverter
	case reflect.Slice, reflect.Array:
		return sliceConverter
	case reflect.Ptr:
		return newPtrConverter
	default:
		return skipConverter
	}
}

func nilConverter(v reflect.Value, opts convertOpts) interface{} {
	return nil
}

func skipConverter(v reflect.Value, opts convertOpts) interface{} {
	return v.Interface()
}

func newPtrConverter(v reflect.Value, opts convertOpts) interface{} {
	if v.IsNil() {
		return nil
	}
	return convert(v.Elem(), opts)
}

func mapConverter(v reflect.Value, opts convertOpts) interface{} {
	if v.IsNil() {
		return nil
	}
	jm := JSONMap{}
	keys := v.MapKeys()
	for _, k := range keys {
		mk := resolveMapKey(k)
		if len(mk) > 0 {
			mv := v.MapIndex(k)
			if opts.recursive {
				jm[mk] = convert(mv, convertOpts{
					recursive:  opts.recursive,
					decorators: opts.decorators,
					path:       fmt.Sprintf("%s.%s", opts.path, mk),
				})
			} else {
				jm[mk] = mv.Interface()
			}
		}
	}
	return jm
}

func resolveMapKey(v reflect.Value) string {
	if v.Kind() == reflect.String {
		return v.String()
	}
	if tm, ok := v.Interface().(encoding.TextMarshaler); ok {
		buf, _ := tm.MarshalText()
		return string(buf)
	}
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	}

	return ""
}

func sliceConverter(v reflect.Value, opts convertOpts) interface{} {
	if v.IsNil() {
		return nil
	}
	if v.Len() == 0 {
		return v.Slice(0, 0).Interface()
	}

	jms := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		jmd := convert(v.Index(i), opts)
		if jm, ok := jmd.(JSONMap); ok {
			if len(opts.decorators) > 0 {
				for _, d := range opts.decorators {
					d(opts.path, jm)
				}
			}
		}
		jms[i] = jmd
	}

	return jms
}

/*func boolConverter(v reflect.Value, opts convertOpts) interface{} {
	return v.Bool()
}

func intConverter(v reflect.Value, opts convertOpts) interface{} {
	return v.Int()
}

func uintConverter(v reflect.Value, opts convertOpts) interface{} {
	return v.Uint()
}

func floatConverter(v reflect.Value, opts convertOpts) interface{} {
	return v.Float()
}

func stringConverter(v reflect.Value, opts convertOpts) interface{} {
	return v.String()
}*/

func interfaceConverter(v reflect.Value, opts convertOpts) interface{} {
	if v.IsNil() {
		return nil
	}

	return convert(v.Elem(), opts)
}

func newStructConverter(t reflect.Type) jmConverter {
	se := structConverter{fields: cachedTypeFields(t)}
	return se.converter
}

type structConverter struct {
	fields []field
}

func (se structConverter) converter(v reflect.Value, opts convertOpts) interface{} {
	jm := JSONMap{}
FieldLoop:
	for i := range se.fields {
		f := &se.fields[i]

		// Find the nested struct field by following f.index.
		fv := v
		for _, i := range f.index {
			if fv.Kind() == reflect.Ptr {
				if fv.IsNil() {
					continue FieldLoop
				}
				fv = fv.Elem()
			}
			fv = fv.Field(i)
		}

		if f.omitEmpty && isEmptyValue(fv) {
			continue
		}
		if opts.recursive {
			jm[f.name] = f.converter(fv, convertOpts{
				recursive:  opts.recursive,
				decorators: opts.decorators,
				path:       fmt.Sprintf("%s.%s", opts.path, f.name),
			})
		} else {
			jm[f.name] = fv.Interface()
		}
	}

	if len(opts.decorators) > 0 {
		for _, d := range opts.decorators {
			d(opts.path, jm)
		}
	}
	return jm
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
