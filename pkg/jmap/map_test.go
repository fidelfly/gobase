package jmap

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestMarshal(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"test",
			args{
				struct {
					A time.Time
					B string
				}{
					time.Now(),
					"testing B",
				},
			},
		},
		{"test2",
			args{
				struct {
					TestS
					A string
				}{
					TestS{
						"Struct A",
						"Struct B",
					},
					"testing A",
				},
			},
		},
		{"test3",
			args{
				struct {
					A string
					B *TestS
				}{
					"testing A",
					&TestS{
						"Struct A",
						"Struct B",
					},
				},
			},
		},
		{"test4",
			args{
				[]struct {
					A string
					B *TestS
				}{{
					"testing A",
					&TestS{
						"Struct A",
						"Struct B",
					},
				},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := jsonStr(tt.args.v)
			got := jsonStr(Marshal(tt.args.v, false))
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Marshal() = %v, want %v", got, want)
			} else {
				t.Logf("json = %s", got)
			}
		})
	}
}

func jmDecorator(path string, jm JSONMap) {
	jm["decorator"] = "Testing decorator"
}

type TestS struct {
	TestSA string
	TestSB string
}

func jsonStr(v interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	data, _ := json.Marshal(v)
	_ = json.Unmarshal(data, &m)
	return m
}
