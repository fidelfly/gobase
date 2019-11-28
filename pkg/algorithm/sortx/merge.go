package sortx

import "reflect"

func MergeSort(sliceOrAddressableArray interface{}, comp Comparator) {
	if ok, sliceVal, _ := isArrayOrSlice(sliceOrAddressableArray); ok {
		aLen := sliceVal.Len()
		if aLen > 0 {
			_mergeSortWithRange(sliceVal, comp, 0, aLen-1)
		}
	}
}

func _mergeSortWithRange(sliceVal reflect.Value, comp Comparator, start int, end int) {
	if start >= end {
		return
	}
	var mid = (start + end) / 2
	_mergeSortWithRange(sliceVal, comp, start, mid)
	_mergeSortWithRange(sliceVal, comp, mid+1, end)
	_mergeSortResult(sliceVal, comp, start, mid, end)
}

func _mergeSortResult(sliceVal reflect.Value, comp Comparator, start int, mid int, end int) {
	aLen := end - start + 1
	temp := reflect.MakeSlice(sliceVal.Type(), aLen, aLen)
	var k, i, j int = 0, start, mid + 1
	for ; k < aLen; k++ {
		if i > mid || j > end {
			break
		}
		if comp(sliceVal.Index(i).Interface(), sliceVal.Index(j).Interface()) > 0 {
			temp.Index(k).Set(sliceVal.Index(j))
			j++
		} else {
			temp.Index(k).Set(sliceVal.Index(i))
			i++
		}
	}

	if i <= mid {
		for ; i <= mid; i++ {
			temp.Index(k).Set(sliceVal.Index(i))
			k++
		}
	}
	if j <= end {
		for ; j <= end; j++ {
			temp.Index(k).Set(sliceVal.Index(j))
			k++
		}
	}

	for k = 0; k < aLen; k++ {
		sliceVal.Index(start + k).Set(temp.Index(k))
	}

}
