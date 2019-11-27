package sortx

import (
	"reflect"
	"testing"
)

func TestQuickSort(t *testing.T) {
	for _, tt := range sortCases {
		t.Run(tt.name, func(t *testing.T) {
			QuickSort(tt.args.target, tt.args.comp)
			if !reflect.DeepEqual(tt.args.target, tt.exp) {
				t.Errorf("Expect : %v, Got : %v", tt.exp, tt.args.target)
			} else {
				t.Logf("Result : %v", tt.args.target)
			}
		})
	}
}
