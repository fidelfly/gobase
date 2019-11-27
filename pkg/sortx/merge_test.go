package sortx

import (
	"reflect"
	"testing"
)

func TestMergeSort(t *testing.T) {
	for _, tt := range sortCases {
		t.Run(tt.name, func(t *testing.T) {
			MergeSort(tt.args.target, tt.args.comp)
			if !reflect.DeepEqual(tt.args.target, tt.exp) {
				t.Errorf("Expect : %v, Got : %v", tt.exp, tt.args.target)
			} else {
				t.Logf("Result : %v", tt.args.target)
			}
		})
	}
}

func Test_MergeSortPtr(t *testing.T) {
	val := []*intStruct{{4}, {5}, {7}, {3}, {3}, {1}, {2}}
	t.Logf("val = %v\n", val)
	val2 := val[:]
	MergeSort(val, intStructPtrComparator)
	t.Logf("val = %v\n", val)
	t.Logf("val2 = %v\n", val2)
}
