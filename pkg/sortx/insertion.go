package sortx

import "reflect"

/*
 * only slice or addressable array is accepted as target
 */
func InsertionSort(sliceOrAddressableArray interface{}, comp Comparator) {
	if ok, sliceVal, swap := isArrayOrSlice(sliceOrAddressableArray); ok {
		aLen := sliceVal.Len()
		for i := 1; i < aLen; i++ {
			val := sliceVal.Index(i).Interface()
			var j = i - 1
			for ; j >= 0; j-- {
				if comp(sliceVal.Index(j).Interface(), val) > 0 {
					swap(j, j+1)
				} else {
					break
				}
			}
			if j+1 != i {
				sliceVal.Index(j + 1).Set(reflect.ValueOf(val))
			}
		}
	}
}
