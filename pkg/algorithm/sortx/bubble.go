package sortx

type Comparator func(a, b interface{}) int

/*
 * only slice or addressable array is accepted as target
 */
func BubbleSort(sliceOrAddressableArray interface{}, comp Comparator) {
	if ok, sliceVal, swap := isArrayOrSlice(sliceOrAddressableArray); ok {
		aLen := sliceVal.Len()
		for i := 0; i < aLen; i++ {
			doSwap := false
			for j := 0; j < aLen-i-1; j++ {
				if comp(sliceVal.Index(j).Interface(), sliceVal.Index(j+1).Interface()) > 0 {
					swap(j, j+1)
					doSwap = true
				}
			}

			if !doSwap {
				return
			}
		}
	}
}
