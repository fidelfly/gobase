package sortx

import "reflect"

func QuickSort(sliceOrAddressableArray interface{}, comp Comparator) {
	if ok, sliceVal, swap := isArrayOrSlice(sliceOrAddressableArray); ok {
		aLen := sliceVal.Len()
		if aLen > 0 {
			_quickSortWithRange(sliceVal, comp, swap, 0, aLen-1)
		}
	}
}

func _quickSortWithRange(sliceVal reflect.Value, comp Comparator, swap func(i, j int), start int, end int) {
	if start >= end {
		return
	}
	q := _quickSortPartition(sliceVal, comp, swap, start, end)
	_quickSortWithRange(sliceVal, comp, swap, start, q-1)
	_quickSortWithRange(sliceVal, comp, swap, q+1, end)
}

func _quickSortPartition(sliceVal reflect.Value, comp Comparator, swap func(i, j int), start int, end int) int {
	var i = start
	valInterface := sliceVal.Index(end).Interface()
	for j := start; j < end; j++ {
		if comp(sliceVal.Index(j).Interface(), valInterface) < 0 {
			swap(i, j)
			i++
		}
	}
	swap(i, end)
	return i
}
