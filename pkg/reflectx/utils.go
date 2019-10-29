package reflectx

import "reflect"

func CopyFields(target interface{}, source interface{}) []string {
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
