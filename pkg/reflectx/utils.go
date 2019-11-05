package reflectx

import "reflect"

func CopyAllFields(target interface{}, source interface{}) []string {
	tv := reflect.Indirect(reflect.ValueOf(target))
	if tv.Kind() != reflect.Struct {
		return nil
	}
	tt := tv.Type()

	sv := reflect.Indirect(reflect.ValueOf(source))
	if sv.Kind() != reflect.Struct {
		return nil
	}
	st := sv.Type()

	var copyFields []string
	for i := 0; i < st.NumField(); i++ {
		sf := st.Field(i)
		fn := sf.Name
		if tf, find := tt.FieldByName(fn); find {
			if tf.Type == sf.Type {
				tfv := tv.FieldByName(fn)
				if tfv.CanSet() {
					tfv.Set(sv.Field(i))
					copyFields = append(copyFields, fn)
				}
			}
		}
	}
	return copyFields
}

func CopyFields(target interface{}, source interface{}, fields ...string) []string {
	if len(fields) == 0 {
		return CopyAllFields(target, source)
	}
	tv := reflect.Indirect(reflect.ValueOf(target))
	if tv.Kind() != reflect.Struct {
		return nil
	}
	tt := tv.Type()

	sv := reflect.Indirect(reflect.ValueOf(source))
	if sv.Kind() != reflect.Struct {
		return nil
	}
	st := sv.Type()

	var copyFields []string
	for _, field := range fields {
		if sf, find := st.FieldByName(field); find {
			if tf, find := tt.FieldByName(field); find {
				if tf.Type == sf.Type {
					tfv := tv.FieldByName(field)
					if tfv.CanSet() {
						tfv.Set(sv.FieldByName(field))
						copyFields = append(copyFields, field)
					}
				}
			}
		}
	}
	return copyFields
}

type FV struct {
	Field string
	Value interface{}
}

func SetField(target interface{}, fvs ...FV) []string {
	tv := reflect.Indirect(reflect.ValueOf(target))
	if tv.Kind() != reflect.Struct {
		return nil
	}
	tt := tv.Type()
	var fields []string
	for _, fv := range fvs {
		if tf, find := tt.FieldByName(fv.Field); find {
			st := reflect.TypeOf(fv.Value)
			if st == tf.Type {
				tfv := tv.FieldByName(fv.Field)
				if tfv.CanSet() {
					tfv.Set(reflect.ValueOf(fv.Value))
					fields = append(fields, fv.Field)
				}
			}
		}
	}
	return fields
}

func GetField(target interface{}, field string) interface{} {
	v := reflect.ValueOf(target)
	if v.IsValid() == false {
		return nil
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}
	t := v.Type()
	if _, find := t.FieldByName(field); find {
		return v.FieldByName(field).Interface()
	}
	return nil
}

func IsValueNil(v interface{}) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}
