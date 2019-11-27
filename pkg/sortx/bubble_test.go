package sortx

import (
	"reflect"
	"testing"
)

type sortArgs struct {
	target interface{}
	comp   Comparator
}

type testCase struct {
	name string
	args sortArgs
	exp  interface{}
}

type intStruct struct {
	Val int
}

var intComparator = func(a, b interface{}) int {
	return a.(int) - b.(int)
}

var intStructComparator = func(a, b interface{}) int {
	return a.(intStruct).Val - b.(intStruct).Val
}

var intStructPtrComparator = func(a, b interface{}) int {
	return a.(*intStruct).Val - b.(*intStruct).Val
}

var sortCases = []testCase{
	// TODO: Add test cases.
	{
		name: "addressable array",
		args: sortArgs{
			target: &[6]int{4, 5, 3, 9, 7, 1},
			comp:   intComparator,
		},
		exp: &[6]int{1, 3, 4, 5, 7, 9},
	},
	{
		name: "slice",
		args: sortArgs{
			target: []int{4, 5, 3, 9, 7, 1},
			comp:   intComparator,
		},
		exp: []int{1, 3, 4, 5, 7, 9},
	},
	{
		name: "struct",
		args: sortArgs{
			target: []intStruct{{4}, {5}, {7}, {3}, {3}, {1}, {2}},
			comp:   intStructComparator,
		},
		exp: []intStruct{{1}, {2}, {3}, {3}, {4}, {5}, {7}},
	},
}

func TestBubbleSort(t *testing.T) {
	for _, tt := range sortCases {
		t.Run(tt.name, func(t *testing.T) {
			BubbleSort(tt.args.target, tt.args.comp)
			if !reflect.DeepEqual(tt.args.target, tt.exp) {
				t.Errorf("Expect : %v, Got : %v", tt.exp, tt.args.target)
			} else {
				t.Logf("Result : %v", tt.args.target)
			}
		})
	}
}

func Test_BubbleSortPtr(t *testing.T) {
	val := []*intStruct{{4}, {5}, {7}, {3}, {3}, {1}, {2}}
	t.Logf("val = %v\n", val)
	val2 := val[:]
	BubbleSort(val, intStructPtrComparator)

	t.Logf("val = %v\n", val)
	t.Logf("val2 = %v\n", val2)
}
