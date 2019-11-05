package reflectx

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	A int
	B int
	s int
}

type TestStruct2 struct {
	A int
	C int
}
type TestStruct3 struct {
	A string
	B int
}

func TestCopyFields(t *testing.T) {
	type args struct {
		target interface{}
		source interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test with same struct",
			args: args{
				target: &TestStruct{
					A: 1,
					B: 2,
					s: 3,
				},
				source: TestStruct{
					A: 11,
					B: 12,
					s: 13,
				},
			},
			want: []string{"A", "B"},
		},
		{
			name: "Test with pointer which has the same struct",
			args: args{
				target: &TestStruct{
					A: 1,
					B: 2,
				},
				source: &TestStruct{
					A: 11,
					B: 12,
				},
			},
			want: []string{"A", "B"},
		},
		{
			name: "Test with different struct",
			args: args{
				target: &TestStruct{
					A: 1,
					B: 2,
				},
				source: TestStruct2{
					A: 11,
					C: 12,
				},
			},
			want: []string{"A"},
		},
		{
			name: "Test with different struct",
			args: args{
				target: &TestStruct{
					A: 1,
					B: 2,
				},
				source: TestStruct3{
					A: "ABC",
					B: 12,
				},
			},
			want: []string{"B"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CopyFields(tt.args.target, tt.args.source); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
			t.Log(tt.args.target)
		})
	}
}

func TestSetField(t *testing.T) {
	ts := []TestStruct{TestStruct{
		A: 1,
		B: 2,
		s: 3,
	}}

	sliceValue := reflect.Indirect(reflect.ValueOf(ts))
	if sliceValue.Kind() == reflect.Slice {
		size := sliceValue.Len()
		if size > 0 {
			for i := 0; i < size; i++ {
				sv := sliceValue.Index(i)
				if sv.Kind() == reflect.Ptr {
					SetField(sv.Interface(), FV{"A", 10})
				} else if sv.CanAddr() {
					SetField(sv.Addr().Interface(), FV{"A", 11})
				}

			}
		}
	}

	for _, v := range ts {
		t.Log(v)
	}

}

type MyType []string

func TestIsValueNil(t *testing.T) {
	type args struct {
		v interface{}
	}
	var v0 *int64
	var f0 func()
	var i0 int64
	var s0 TestStruct
	var mt MyType
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Nil *int64",
			args{v0},
			true,
		},
		{
			name: "Struct",
			args: args{TestStruct{1, 2, 3}},
			want: false,
		},
		{
			name: "UnInit Struct",
			args: args{s0},
			want: false,
		},
		{
			name: "Nil Struct Pointer",
			args: args{returnNilStruct()},
			want: true,
		},
		{
			name: "string",
			args: args{"String"},
			want: false,
		},
		{
			name: "Nil Function",
			args: args{f0},
			want: true,
		},
		{
			name: "Function",
			args: args{func() {}},
			want: false,
		},
		{
			name: "int",
			args: args{i0},
			want: false,
		},
		{
			name: "Alias Type",
			args: args{mt},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Direct compare value == nil : %t", tt.args.v == nil)
			if got := IsValueNil(tt.args.v); got != tt.want {
				t.Errorf("IsValueNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func returnNilStruct() interface{} {
	var a *TestStruct

	return a
}
