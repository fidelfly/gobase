package sortx

/*
 * only slice or addressable array is accepted as target
 */
func SelectionSort(sliceOrAddressableArray interface{}, comp Comparator) {
	if ok, sliceVal, swap := isArrayOrSlice(sliceOrAddressableArray); ok {
		aLen := sliceVal.Len()
		for i := 0; i < aLen; i++ {
			minVal := sliceVal.Index(i).Interface()
			minIndex := i
			for j := i + 1; j < aLen; j++ {
				val := sliceVal.Index(j).Interface()
				if comp(val, minVal) < 0 {
					minVal = val
					minIndex = j
				}
			}
			if minIndex != i {
				swap(i, minIndex)
			}
		}
	}
}
