package sortx

import "reflect"

func isArrayOrSlice(target interface{}) (bool, reflect.Value, func(i, j int)) {
	val := reflect.Indirect(reflect.ValueOf(target))
	switch val.Kind() {
	case reflect.Array:
		if !val.CanAddr() {
			return false, val, nil
		}
		sliceVal := val.Slice(0, val.Len())
		return true, sliceVal, reflect.Swapper(sliceVal.Interface())
	case reflect.Slice:
		return true, val, reflect.Swapper(target)
	default:
		return false, val, nil
	}
}
