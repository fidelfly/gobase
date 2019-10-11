package randx

import "testing"

func TestCharPool_Get(t *testing.T) {
	type args struct {
		len int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			"test1",
			args{16},
		},
		{
			"test2",
			args{8},
		},
		{
			"test2",
			args{6},
		},
		{
			"test2",
			args{8},
		},
		{
			"test2",
			args{4},
		},
		{
			"test3",
			args{12},
		},
	}
	cp := NewCharPool(
		NewCharSource("abcdefg中文测试ijklmnopqrstuvwxyz123456789", 8),
		NewCharSource("+_)(*&", 2),
		NewCharSource("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 2),
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Get(%d) = %s", tt.args.len, cp.Get(tt.args.len))
		})
	}
}
